package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/kk-digital/kcg-devops-gogs-mirror/pkg/client"
	"github.com/spf13/cobra"
)

var (
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
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Add code to update Gogs repos
			log.Println("Updating Gogs repositories...")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&githubAccessToken, "github-token", "t", "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", "GitHub access token")
	rootCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-url", "u", "localhost:10880", "Gogs base URL")
	rootCmd.PersistentFlags().StringVarP(&gogsSSHURL, "gogs-ssh-url", "s", "localhost:10022", "Gogs ssh URL")
	rootCmd.PersistentFlags().StringVarP(&gogsUserName, "gogs-user-name", "n", "my-name", "your Gogs user name")
	rootCmd.PersistentFlags().StringVarP(&gogsAccessToken, "gogs-token", "g", "77cae12a2134d6e6ad8da5262a90502a412d7c03", "Gogs base URL")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", runtime.NumCPU(), "Speed up the command")

	cloneCmd.PersistentFlags().StringVarP(&orgName, "org-name", "o", "demo-33383080", "grabs all repos from an organization")

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
	gogsClient := client.NewGogsClient(gogsBaseURL, gogsSSHURL, gogsUserName, gogsAccessToken)
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

		repoName, cloneURL := *repo.Name, *repo.CloneURL
		// 3.1. create gogs org repo if not exists
		gogsRepo, err := gogsClient.GetOrgRepo(orgName, repoName)
		if err != nil {
			log.Fatal(err)
		}
		if len(gogsRepo) != 0 {
			log.Printf("Repository %s already exists, skipped, cost: %s\n", *repo.FullName, time.Since(now))
			continue
		}

		// 3.2. clone repo to gogs org repo
		if err = gogsClient.CloneRepoToGogs(orgName, repoName, cloneURL); err != nil {
			log.Fatal(err)
		}
		log.Printf("Repository %s, cost: %s\n", *repo.FullName, time.Since(now))
	}

	log.Printf("Cloned successfully, total cost: %s\n", time.Since(cloneNow))
}
