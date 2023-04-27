package cmd

import (
	"context"
	"log"
	"time"

	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/client"
	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/pretty"
	"github.com/spf13/cobra"
)

var listOrgRepoCmd = &cobra.Command{
	Use:   "list-org-repo",
	Short: "Get a list of github repositories in an organization",
	Run:   listOrgRepo,
}

func init() {
	listOrgRepoCmd.PersistentFlags().StringVarP(&githubAccessToken, "github-token", "t", "", "GitHub access token")
	listOrgRepoCmd.PersistentFlags().StringVarP(&orgName, "org-name", "o", "", "grabs all repos from an organization")

	listOrgRepoCmd.MarkPersistentFlagRequired("github-token")
	listOrgRepoCmd.MarkPersistentFlagRequired("org-name")

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
		log.Println(pretty.JSON(repo))
	}

	log.Printf("Successfully listed repositories, total cost: %s\n", time.Since(now))
}
