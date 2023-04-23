package cmd

import (
	"context"
	"log"
	"time"

	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/client"
	"github.com/spf13/cobra"
)

var listOrgRepoCmd = &cobra.Command{
	Use:   "list-org-repo",
	Short: "Get a list of github repositories in an organization",
	Run:   listOrgRepo,
}

func init() {
	rootCmd.AddCommand(listOrgRepoCmd)
}

func listOrgRepo(cmd *cobra.Command, args []string) {
	now := time.Now()
	log.Println("Listing GitHub repositories...")

	ctx := context.Background()
	githubClient := client.NewGithubClient(ctx, githubAccessToken)
	allRepos, err := githubClient.ListOrgRepos(ctx, orgName)
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range allRepos {
		log.Println(repo)
	}

	log.Printf("Successfully listed repositories, total cost: %s\n", time.Since(now))
}
