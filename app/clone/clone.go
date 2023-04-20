package clone

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/google/go-github/v51/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// TODO: Add code to clone GitHub repos to Gogs
func Command(cmd *cobra.Command, args []string) {
	fmt.Println("Cloning GitHub repositories to Gogs...")
	fmt.Println(args)

	// githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	// gogsBaseURL := os.Getenv("GOGS_BASE_URL")
	githubAccessToken := "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"
	gogsBaseURL := "/home/rabbit/gogs/git/gogs-repositories/"

	if githubAccessToken == "" || gogsBaseURL == "" {
		log.Fatal("Please set GITHUB_ACCESS_TOKEN and GOGS_BASE_URL environment variables")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Replace "myorg" with your GitHub organization name
	// org := "asdv23"

	// opt := &github.RepositoryListByOrgOptions{
	// 	ListOptions: github.ListOptions{PerPage: 10},
	// }
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	for {
		repos, resp, err := client.Repositories.List(ctx, "", nil)
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range repos {
			cloneURL := fmt.Sprintf("%s%s.git", gogsBaseURL, *repo.FullName)
			fmt.Printf("Cloning %s to %s...\n", *repo.CloneURL, cloneURL)
			cmd := exec.Command("git", "clone", "--mirror", *repo.CloneURL, cloneURL)
			err := cmd.Run()
			if err != nil {
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
