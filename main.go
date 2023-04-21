package main

import (
	"log"
	"os"

	"github.com/kk-digital/kcg-devops-gogs-mirror/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
