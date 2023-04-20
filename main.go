package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/go-github/v51/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	githubAccessToken string
	gogsBaseURL       string
	gogsSSHURL        string
	gogsOrgName       string
	gogsAccessToken   string
	workers           int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gogs-helper",
		Short: "A helper tool to clone and update repositories between GitHub and Gogs",
	}

	cloneCmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone all repos from GitHub organization to Gogs",
		Run:   clone,
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update all existing repos in Gogs",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Add code to update Gogs repos
			fmt.Println("Updating Gogs repositories...")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&githubAccessToken, "github-token", "t", "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", "GitHub access token")
	rootCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-url", "u", "http://localhost:10880", "Gogs base URL")
	rootCmd.PersistentFlags().StringVarP(&gogsSSHURL, "gogs-ssh-url", "s", "ssh://git@localhost:10022", "Gogs ssh URL")
	rootCmd.PersistentFlags().StringVarP(&gogsOrgName, "gogs-org-name", "o", "my-name", "your Gogs organization name")
	rootCmd.PersistentFlags().StringVarP(&gogsAccessToken, "gogs-token", "g", "77cae12a2134d6e6ad8da5262a90502a412d7c03", "Gogs base URL")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", runtime.NumCPU(), "Gogs base URL")
	// rootCmd.MarkPersistentFlagRequired("github-token")
	// rootCmd.MarkPersistentFlagRequired("gogs-url")

	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(updateCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func clone(cmd *cobra.Command, args []string) {
	fmt.Println("Cloning GitHub repositories to Gogs...")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		repos, resp, err := client.Repositories.List(ctx, "", nil)
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range repos {
			fmt.Printf("Cloning %s to gogs...\n", *repo.CloneURL)
			if err = cloneRepoToGogs(repo); err != nil {
				log.Printf("Failed to clone %s: %v\n", *repo.FullName, err)
			} else {
				fmt.Printf("Successfully cloned %s\n", *repo.FullName)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}

func cloneRepoToGogs(repo *github.Repository) error {
	repoName := *repo.Name

	// First, create the repository in Gogs using the Gogs API
	err := createGogsRepo(repoName)
	if err != nil {
		return fmt.Errorf("failed to create Gogs repository: %w", err)
	}

	// Use the git command to clone the GitHub repository and then push to the Gogs repository
	cmd := exec.Command("git", "clone", "--mirror", *repo.CloneURL)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone GitHub repository: %w", err)
	}

	// Change to the cloned repository's directory
	err = os.Chdir(repoName + ".git")
	if err != nil {
		return fmt.Errorf("failed to change to the cloned repository directory: %w", err)
	}

	// Construct the Gogs repository URL with the token for authentication
	gogsRepoURL := fmt.Sprintf("%s/%s/%s.git", gogsSSHURL, gogsOrgName, repoName)

	// Add the Gogs remote
	cmd = exec.Command("git", "remote", "add", "gogs", gogsRepoURL)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to add Gogs remote: %w", err)
	}

	// Push the cloned repository to the Gogs remote
	cmd = exec.Command("git", "push", "--mirror", "gogs")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to push to Gogs repository: %w", err)
	}

	// Change back to the original directory
	if err = os.Chdir(".."); err != nil {
		return fmt.Errorf("failed to change back to the original directory: %w", err)
	}

	// Remove the cloned repository
	cmd = exec.Command("rm", "-rf", repoName+".git")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove the cloned repository: %w", err)
	}

	return nil
}

type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private,omitempty"`
}

func createGogsRepo(repoName string) error {
	createRepoRequest := &CreateRepoRequest{
		Name:        repoName,
		Description: "Cloned from GitHub",
		Private:     true,
	}

	jsonData, err := json.Marshal(createRepoRequest)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	// https://github.com/gogs/docs-api/tree/master/Repositories#create
	// gogsRepoURL := fmt.Sprintf("%s/api/v1/org/%s/repos", gogsBaseURL, orgName)
	gogsRepoURL := fmt.Sprintf("%s/api/v1/user/repos", gogsBaseURL)
	req, err := http.NewRequest(http.MethodPost, gogsRepoURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+gogsAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error executing request: %v", err)
		}
		fmt.Println("data:", string(data))
		fmt.Println("resp:", resp)
		return fmt.Errorf("failed to create repository. Status code: %d", resp.StatusCode)
	}

	fmt.Printf("Repository %s created successfully.\n", repoName)
	return nil
}
