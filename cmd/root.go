package cmd

import (
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

	// gogs auth
	user      string
	pass      string
	tokenName string

	// gogs ssh
	title string
	key   string
)

var rootCmd = &cobra.Command{
	Use:   "gogs-helper",
	Short: "A helper tool to clone and update repositories between GitHub and Gogs",
}

func Execute() error {
	return rootCmd.Execute()
}
