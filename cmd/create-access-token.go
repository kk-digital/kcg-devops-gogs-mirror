package cmd

import (
	"log"
	"time"

	"github.com/gogs/go-gogs-client"
	"github.com/spf13/cobra"
)

var createAccessTokenCmd = &cobra.Command{
	Use:   "create-access-token",
	Short: "Create an access token with a user and password",
	Run:   createAccessToken,
}

func init() {
	createAccessTokenCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-base-url", "b", "", "Gogs base URL, e.g. http://localhost:10880")
	createAccessTokenCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "System user, e.g. admin")
	createAccessTokenCmd.PersistentFlags().StringVarP(&pass, "pass", "p", "", "System user password, e.g. password")
	createAccessTokenCmd.PersistentFlags().StringVarP(&tokenName, "token-name", "n", "", "Access token name, e.g. script_token")

	createAccessTokenCmd.MarkPersistentFlagRequired("gogs-base-url")
	createAccessTokenCmd.MarkPersistentFlagRequired("user")
	createAccessTokenCmd.MarkPersistentFlagRequired("pass")
	createAccessTokenCmd.MarkPersistentFlagRequired("token-name")

	rootCmd.AddCommand(createAccessTokenCmd)
}

func createAccessToken(cmd *cobra.Command, args []string) {
	now := time.Now()
	log.Println("creating access token with gogs...")

	client := gogs.NewClient(gogsBaseURL, "")
	opt := gogs.CreateAccessTokenOption{
		Name: tokenName,
	}
	accessToken, err := client.CreateAccessToken(user, pass, opt)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user %s created access token, Name: %s, Sha1: %s\n", user, accessToken.Name, accessToken.Sha1)

	log.Printf("Successfully created acccess token, total cost: %s\n", time.Since(now))
}
