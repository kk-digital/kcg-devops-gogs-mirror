package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/client"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone all repos from GitHub organization to Gogs",
	Run:   clone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func clone(cmd *cobra.Command, args []string) {
	cloneNow := time.Now()
	log.Println("Cloning GitHub repositories to Gogs...")

	// 1. create gogs org if not exists
	gogsClient := client.NewGogsClient(gogsBaseURL, user, gogsAccessToken)
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
	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Mkdir repos
	repoDir := filepath.Join(workdir, repoName+".git")
	if err := cloneLocal_(repoName, cloneURL); err != nil {
		return fmt.Errorf("failed to clone repositories into local directory: %w", err)
	}

	// Change to the cloned repository's directory
	err = os.Chdir(repoDir)
	if err != nil {
		return fmt.Errorf("failed to change to the cloned repository directory: %w", err)
	}

	// Construct the Gogs repository URL with the token for authentication
	gogsSSHURL_ := "ssh://git@" + gogsSSHURL
	gogsRepoURL := fmt.Sprintf("%s/%s/%s.git", gogsSSHURL_, orgName, repoName)

	// Add the Gogs remote
	cmd := exec.Command("git", "remote", "add", "gogs", gogsRepoURL)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to add Gogs remote: %w", err)
	}

	// Push the cloned repository to the Gogs remote
	cmd = exec.Command("git", "push", "--mirror", "gogs")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to push to Gogs repository %s: %w", repoName, err)
	}

	// Change back to the original directory
	if err = os.Chdir(dir); err != nil {
		return fmt.Errorf("failed to change back to the original directory: %w", err)
	}

	return nil
}
