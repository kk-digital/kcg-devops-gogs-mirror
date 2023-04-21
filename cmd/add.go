package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/client"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add all repositories from a local directory to Gogs",
	Run:   add,
}

func init() {
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

	gogsClient := client.NewGogsClient(gogsBaseURL, gogsUserName, gogsAccessToken)

	err = filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk %s: %w", workdir, err)
		}
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			now := time.Now()

			repoName := strings.TrimSuffix(filepath.Base(path), ".git")
			// First, create the repository in Gogs using the Gogs API
			if err = gogsClient.CreateRepoInOrg(orgName, repoName); err != nil {
				return fmt.Errorf("failed to create repository to gogs: %w", err)
			}

			// Change to the cloned repository's directory
			err = os.Chdir(path)
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
