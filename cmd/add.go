package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogs/go-gogs-client"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add all repositories from a local directory to Gogs",
	Run:   add,
}

func init() {
	addCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-base-url", "b", "", "Gogs base URL, e.g. http://localhost:10880")
	addCmd.PersistentFlags().StringVarP(&gogsSSHURL, "gogs-ssh-url", "s", "", "Gogs ssh URL, e.g. ssh://git@localhost:10022")
	addCmd.PersistentFlags().StringVarP(&gogsAccessToken, "gogs-token", "t", "", "Gogs access token, e.g. 221a1527091612fade38d265742b84c40ab17de1")
	addCmd.PersistentFlags().StringVarP(&orgName, "org-name", "o", "", "Add all repos to an organization")
	addCmd.PersistentFlags().StringVarP(&workdir, "workdir", "d", "", "The working directory will store all the repository of github")

	addCmd.MarkPersistentFlagRequired("gogs-base-url")
	addCmd.MarkPersistentFlagRequired("gogs-ssh-url")
	addCmd.MarkPersistentFlagRequired("gogs-token")
	addCmd.MarkPersistentFlagRequired("workdir")
	rootCmd.AddCommand(addCmd)
}

func add(cmd *cobra.Command, args []string) {
	now := time.Now()
	log.Println("Adding repositories to gogs...")

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	client := gogs.NewClient(gogsBaseURL, gogsAccessToken)
	// 1. create gogs org if not exists
	if _, err = client.GetOrg(orgName); err != nil && err.Error() == "404 Not Found" {
		if _, err = client.CreateOrg(gogs.CreateOrgOption{
			UserName:    orgName,
			FullName:    orgName,
			Description: "Cloned organization from GitHub",
		}); err != nil {
			log.Fatal(fmt.Errorf("failed to create organization %s: %w", orgName, err))
		}
		log.Printf("Successfully added organization %s\n", orgName)
	}
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get organization %s: %w", orgName, err))
	}

	err = filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk %s: %w", workdir, err)
		}
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			now := time.Now()

			repoName := strings.TrimSuffix(filepath.Base(path), ".git")
			// First, create the repository in Gogs using the Gogs API
			if _, err = client.CreateOrgRepo(orgName, gogs.CreateRepoOption{
				Name:        repoName,
				Description: "Cloned from GitHub",
				Private:     true,
			}); err != nil {
				return fmt.Errorf("failed to create repository to gogs: %w", err)
			}

			// Change to the cloned repository's directory
			err = os.Chdir(path)
			if err != nil {
				return fmt.Errorf("failed to change to the cloned repository directory: %w", err)
			}

			// Construct the Gogs repository URL with the token for authentication
			gogsRepoURL := fmt.Sprintf("%s/%s/%s.git", gogsSSHURL, orgName, repoName)

			// Add the Gogs remote
			cmd := exec.Command("git", "remote", "add", "gogs", gogsRepoURL)
			if err = cmd.Run(); err != nil {
				fmt.Println("err:", err)
				return fmt.Errorf("failed to add Gogs remote %s: %w", gogsRepoURL, err)
			}

			// Push the updated repository to the Gogs remote
			cmd = exec.Command("git", "push", "--mirror", "gogs")
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("failed to push to Gogs repository %s: %w", path, err)
			}

			// Change back to the original directory
			if err = os.Chdir(dir); err != nil {
				return fmt.Errorf("failed to change back to the original directory: %w", err)
			}

			log.Printf("Adding repository %s to gogs, cost: %s\n", path, time.Since(now))
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error while adding repos: %v", err)
	}

	log.Printf("Successfully added repositories, total cost: %s\n", time.Since(now))
}
