package main

import (
	"fmt"
	"github.com/levin-kitty/go-test/app"
	"log"
	"os"
)

func main() {
	privateKeyPath := os.Getenv("KEY_PATH")
	if privateKeyPath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "KEY_PATH not set\n")
		os.Exit(1)
	}

	serverApiBaseUrl := os.Getenv("SERVER_API_BASE_URL")
	if serverApiBaseUrl == "" {
		_, _ = fmt.Fprintf(os.Stderr, "SERVER_API_BASE_URL not set\n")
		os.Exit(1)
	}

	apple, err := app.NewApp(privateKeyPath, serverApiBaseUrl)
	if err != nil {
		panic(err)
	}
	log.Fatalln(apple.Run("127.0.0.1:8080"))
}
