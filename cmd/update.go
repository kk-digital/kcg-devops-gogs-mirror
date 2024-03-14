package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all existing repos in Gogs",
	Run:   update,
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&workdir, "workdir", "d", "", "The working directory will store all the repository of github")

	rootCmd.AddCommand(updateCmd)
}

func update(cmd *cobra.Command, args []string) {
	updateNow := time.Now()
	log.Println("Updating Gogs repositories from Github...")

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	err = filepath.Walk(workdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk %s: %w", workdir, err)
		}
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			now := time.Now()

			// Update from github
			cmd := exec.Command("git", "-C", path, "remote", "update")
			if err := cmd.Run(); err != nil {
				log.Printf("failed to update %s: %w", path, err)
				return filepath.SkipDir
			}

			// Change to the cloned repository's directory
			err = os.Chdir(path)
			if err != nil {
				return fmt.Errorf("failed to change to the cloned repository directory: %w", err)
			}

			// Push the updated repository to the Gogs remote
			cmd = exec.Command("git", "push", "--mirror", "gogs")
			if err = cmd.Run(); err != nil {
				log.Printf("failed to push to Gogs repository %s: %w", path, err)

				// Change back to the original directory
				if err = os.Chdir(dir); err != nil {
					return fmt.Errorf("failed to change back to the original directory: %w", err)
				}

				return filepath.SkipDir
			}

			// Change back to the original directory
			if err = os.Chdir(dir); err != nil {
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
