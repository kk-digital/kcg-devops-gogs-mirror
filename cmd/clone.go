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
	Use:   "clone-local",
	Short: "Clone all repos from GitHub organization into a local directory",
	Long:  "Clone all repos from GitHub organization into a local directory. Duplicate clone will fail: exit status 128",
	Run:   clone,
}

func init() {
	cloneCmd.PersistentFlags().StringVarP(&githubAccessToken, "github-token", "t", "", "GitHub access token")
	cloneCmd.PersistentFlags().StringVarP(&orgName, "org-name", "o", "", "grabs all repos from an organization")
	cloneCmd.PersistentFlags().StringVarP(&workdir, "workdir", "d", "", "The working directory will store all the repository of github")

	cloneCmd.MarkPersistentFlagRequired("github-token")
	cloneCmd.MarkPersistentFlagRequired("org-name")
	cloneCmd.MarkPersistentFlagRequired("workdir")

	rootCmd.AddCommand(cloneCmd)
}

func clone(cmd *cobra.Command, args []string) {
	now := time.Now()
	log.Println("Cloning GitHub repositories to local directory...")

	ctx := context.Background()
	githubClient := client.NewGithubClient(ctx, githubAccessToken)
	allRepos, err := githubClient.ListOrgRepos(ctx, orgName)
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range allRepos {
		cloneNow := time.Now()
		if err = clone_(*repo.Name, *repo.SSHURL); err != nil {
			log.Fatal(err)
		}

		log.Printf("Cloning repository %s, cost: %s\n", *repo.FullName, time.Since(cloneNow))
	}

	log.Printf("Successfully cloned repositories, total cost: %s\n", time.Since(now))
}

func clone_(repoName, cloneURL string) error {
	// Mkdir repos
	repoDir := filepath.Join(workdir, repoName+".git")
	if err := os.MkdirAll(repoDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to mkdir gogs repositories directory: %w", err)
	}

	// Use the git command to clone the GitHub repository and then push to the Gogs repository
	cmd := exec.Command("git", "clone", "--mirror", cloneURL, repoDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone GitHub repository: %w", err)
	}

	return nil
}
