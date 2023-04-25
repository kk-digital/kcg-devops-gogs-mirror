package cmd

import (
	"log"
	"time"

	"github.com/gogs/go-gogs-client"
	"github.com/spf13/cobra"
)

var createPublicKeyCmd = &cobra.Command{
	Use:   "create-public-key",
	Short: "Create a public key to use git ssh",
	Run:   createPublicKey,
}

func init() {
	createPublicKeyCmd.PersistentFlags().StringVarP(&gogsBaseURL, "gogs-base-url", "b", "", "Gogs base URL, e.g. http://localhost:10880")
	createPublicKeyCmd.PersistentFlags().StringVarP(&gogsAccessToken, "gogs-token", "t", "", "Gogs access token, e.g. 221a1527091612fade38d265742b84c40ab17de1")
	createPublicKeyCmd.PersistentFlags().StringVarP(&title, "title", "s", "", "SSH title, e.g. ssh-rsa")
	createPublicKeyCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "SSH public key, e.g. cat ~/.ssh/id_rsa.pub")

	createPublicKeyCmd.MarkPersistentFlagRequired("gogs-base-url")
	createPublicKeyCmd.MarkPersistentFlagRequired("gogs-token")
	createPublicKeyCmd.MarkPersistentFlagRequired("title")
	createPublicKeyCmd.MarkPersistentFlagRequired("key")

	rootCmd.AddCommand(createPublicKeyCmd)
}

func createPublicKey(cmd *cobra.Command, args []string) {
	now := time.Now()
	log.Println("creating public key with gogs...")

	client := gogs.NewClient(gogsBaseURL, gogsAccessToken)
	publicKey, err := client.CreatePublicKey(gogs.CreateKeyOption{
		Title: title,
		Key:   key,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created public key, Title: %s, URL: %s\n", publicKey.Title, publicKey.URL)

	log.Printf("Successfully created public key, total cost: %s\n", time.Since(now))
}
