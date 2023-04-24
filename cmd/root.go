package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// save repos
	workdir string

	// github
	githubAccessToken string

	// organazition
	orgName string

	// gogs
	gogsBaseURL     string
	gogsSSHURL      string
	gogsAccessToken string

	// concurrency
	workers int

	// gogs pam
	user      string
	pass      string
	tokenName string
)

var rootCmd = &cobra.Command{
	Use:   "gogs-helper",
	Short: "A helper tool to clone and update repositories between GitHub and Gogs",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&workdir, "workdir", "d", "repos", "The working directory will store all the repository of github")
	rootCmd.PersistentFlags().StringVarP(&githubAccessToken, "github-token", "t", "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", "GitHub access token")
	rootCmd.PersistentFlags().StringVarP(&orgName, "org-name", "o", "demo-33383080", "grabs all repos from an organization")
	rootCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-http-url", "b", "localhost:10880", "Gogs base URL")
	rootCmd.PersistentFlags().StringVarP(&gogsSSHURL, "gogs-ssh-url", "s", "localhost:10022", "Gogs ssh URL")
	rootCmd.PersistentFlags().StringVarP(&gogsAccessToken, "gogs-token", "g", "77cae12a2134d6e6ad8da5262a90502a412d7c03", "Gogs base URL")
	rootCmd.PersistentFlags().IntVarP(&workers, "workers", "w", runtime.NumCPU(), "Speed up the command")

	// rootCmd.MarkPersistentFlagRequired("github-token")
	// rootCmd.MarkPersistentFlagRequired("gogs-token")
}

func Execute() error {
	return rootCmd.Execute()
}
