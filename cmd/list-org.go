package cmd

import (
	"context"
	"log"
	"time"

	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/client"
	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/pretty"
	"github.com/spf13/cobra"
)

var listOrgCmd = &cobra.Command{
	Use:   "list-org",
	Short: "Get a list of github organizations",
	Run:   listOrg,
}

func init() {
	listOrgCmd.PersistentFlags().StringVarP(&githubAccessToken, "github-token", "t", "", "GitHub access token")

	listOrgCmd.MarkPersistentFlagRequired("github-token")

	rootCmd.AddCommand(listOrgCmd)
}

func listOrg(cmd *cobra.Command, args []string) {
	now := time.Now()
	log.Println("Listing GitHub organizations...")

	ctx := context.Background()
	githubClient := client.NewGithubClient(ctx, githubAccessToken)
	allOrgs, err := githubClient.ListAllOrgs(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, org := range allOrgs {
		log.Println(pretty.JSON(org))
	}

	log.Printf("Successfully listed organizations, total cost: %s\n", time.Since(now))
}
