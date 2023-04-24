package cmd

import (
	"log"
	"time"

	"github.com/gogs/go-gogs-client"
	"github.com/spf13/cobra"
)

var createAccessTokenCmd = &cobra.Command{
	Use:   "create-access-token",
	Short: "sign up a user with Gogs",
	Run:   createAccessToken,
}

func init() {
	createAccessTokenCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-http-url", "b", "localhost:10880", "Gogs base URL")
	createAccessTokenCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "System user")
	createAccessTokenCmd.PersistentFlags().StringVarP(&pass, "pass", "p", "", "System user password")
	createAccessTokenCmd.PersistentFlags().StringVarP(&tokenName, "token-name", "n", "script-token", "Access token name")

	createAccessTokenCmd.MarkPersistentFlagRequired("user")
	createAccessTokenCmd.MarkPersistentFlagRequired("pass")

	rootCmd.AddCommand(createAccessTokenCmd)
}

func createAccessToken(cmd *cobra.Command, args []string) {
	now := time.Now()
	log.Println("Getting access token with gogs...")

	client := gogs.NewClient("http://"+gogsBaseURL, "")
	opt := gogs.CreateAccessTokenOption{
		Name: tokenName,
	}
	accessToken, err := client.CreateAccessToken("rabbit", "rabbit", opt)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user %s created access token, Name: %s, Sha1: %s\n", user, accessToken.Name, accessToken.Sha1)

	log.Printf("Successfully getted acccess token, total cost: %s\n", time.Since(now))
}
