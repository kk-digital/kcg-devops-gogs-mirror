package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/client"
	"github.com/spf13/cobra"
)

var (
	// save repos
	workdir string

	// github
	githubAccessToken string

	// organazition
	orgName string

	// gogs
	gogsBaseURL     string
	gogsSSHURL      string
	gogsUserName    string
	gogsAccessToken string

	// concurrency
	workers int
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
		Run:   update,
	}

	rootCmd.PersistentFlags().StringVarP(&workdir, "workdir", "d", ".", "The working directory will store all the repository of github under its subdirectory repos")
	rootCmd.PersistentFlags().StringVarP(&githubAccessToken, "github-token", "t", "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", "GitHub access token")
	rootCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-url", "u", "localhost:10880", "Gogs base URL")
	rootCmd.PersistentFlags().StringVarP(&gogsSSHURL, "gogs-ssh-url", "s", "localhost:10022", "Gogs ssh URL")
	rootCmd.PersistentFlags().StringVarP(&gogsUserName, "gogs-user-name", "n", "my-name", "your Gogs user name")
	rootCmd.PersistentFlags().StringVarP(&gogsAccessToken, "gogs-token", "g", "77cae12a2134d6e6ad8da5262a90502a412d7c03", "Gogs base URL")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", runtime.NumCPU(), "Speed up the command")

	cloneCmd.PersistentFlags().StringVarP(&orgName, "org-name", "o", "demo-33383080", "grabs all repos from an organization")

	rootCmd.MarkPersistentFlagRequired("github-token")
	rootCmd.MarkPersistentFlagRequired("gogs-token")

	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(updateCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func clone(cmd *cobra.Command, args []string) {
	cloneNow := time.Now()
	log.Println("Cloning GitHub repositories to Gogs...")

	// 1. create gogs org if not exists
	gogsClient := client.NewGogsClient(gogsBaseURL, gogsUserName, gogsAccessToken)
	gogsOrg, err := gogsClient.GetOrg(orgName)
	if err != nil {
		log.Fatal(err)
	}
	if len(gogsOrg) == 0 {
		if err = gogsClient.CreateOrg(orgName); err != nil {
			log.Fatal(err)
		}
	}

	// 2. input org, find all the github repo
	ctx := context.Background()
	githubClient := client.NewGithubClient(ctx, githubAccessToken)
	allRepos, err := githubClient.ListOrgRepos(ctx, orgName)
	if err != nil {
		log.Fatal(err)
	}

	// 3. clone repo to gogs org
	for _, repo := range allRepos {
		now := time.Now()

		repoName, cloneURL := *repo.Name, *repo.SSHURL
		// 3.1. create gogs org repo if not exists
		gogsRepo, err := gogsClient.GetOrgRepo(orgName, repoName)
		if err != nil {
			log.Fatal(err)
		}
		if len(gogsRepo) != 0 {
			log.Printf("Repository %s already exists, skipped, cost: %s\n", *repo.FullName, time.Since(now))
			continue
		}

		// First, create the repository in Gogs using the Gogs API
		if err = gogsClient.CreateRepoInOrg(orgName, repoName); err != nil {
			log.Fatal(err)
		}

		// 3.2. clone repo to gogs org repo
		if err = clone_(repoName, cloneURL); err != nil {
			log.Fatal(err)
		}

		log.Printf("Cloning repository %s, cost: %s\n", *repo.FullName, time.Since(now))
	}

	log.Printf("Successfully cloned, total cost: %s\n", time.Since(cloneNow))
}

func clone_(repoName, cloneURL string) error {
	var err error

	// Mkdir repos
	repoDir := filepath.Join(workdir, "repos", repoName+".git")
	if err := os.MkdirAll(repoDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to mkdir gogs repositories directory: %w", err)
	}

	// Use the git command to clone the GitHub repository and then push to the Gogs repository
	cmd := exec.Command("git", "clone", "--mirror", cloneURL, repoDir)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone GitHub repository: %w", err)
	}

	// Change to the cloned repository's directory
	err = os.Chdir(repoDir)
	if err != nil {
		return fmt.Errorf("failed to change to the cloned repository directory: %w", err)
	}

	// Construct the Gogs repository URL with the token for authentication
	gogsSSHURL := "ssh://git@" + gogsSSHURL
	gogsRepoURL := fmt.Sprintf("%s/%s/%s.git", gogsSSHURL, orgName, repoName)

	// Add the Gogs remote
	cmd = exec.Command("git", "remote", "add", "gogs", gogsRepoURL)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to add Gogs remote: %w", err)
	}

	// Push the cloned repository to the Gogs remote
	cmd = exec.Command("git", "push", "--mirror", "gogs")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to push to Gogs repository %s: %w", repoName, err)
	}

	// Change back to the original directory
	if err = os.Chdir("../.."); err != nil {
		return fmt.Errorf("failed to change back to the original directory: %w", err)
	}

	return nil
}

func update(cmd *cobra.Command, args []string) {
	updateNow := time.Now()
	log.Println("Updating Gogs repositories from Github...")

	reposPath := workdir + "/repos"
	err := filepath.Walk(reposPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk %s: %w", reposPath, err)
		}
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			now := time.Now()

			// Update from github
			cmd := exec.Command("git", "-C", path, "remote", "update")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to update %s: %w", path, err)
			}

			// Change to the cloned repository's directory
			err = os.Chdir(path)
			if err != nil {
				return fmt.Errorf("failed to change to the cloned repository directory: %w", err)
			}

			// Push the updated repository to the Gogs remote
			cmd = exec.Command("git", "push", "--mirror", "gogs")
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("failed to push to Gogs repository %s: %w", path, err)
			}

			// Change back to the original directory
			if err = os.Chdir("../.."); err != nil {
				return fmt.Errorf("failed to change back to the original directory: %w", err)
			}

			log.Printf("Updating repository %s, cost: %s\n", path, time.Since(now))
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error while updating repos: %v", err)
	}

	log.Printf("Successfully updated, total cost: %s\n", time.Since(updateNow))
}
