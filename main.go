package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

type Repository struct {
	FullName string `json:"full_name"`
	Owner    struct {
		DisplayName string `json:"display_name"`
	}
}

func getRepositories() []Repository {
	godotenv.Load()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.bitbucket.org/2.0/repositories", nil)
	if err != nil {
		panic(err)
	}
	bbUser := os.Getenv("BB_USER")
	bbPass := os.Getenv("BB_PASS")
	req.SetBasicAuth(bbUser, bbPass)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response struct {
		Values []Repository
	}
	json.Unmarshal(body, &response)
	return response.Values
}

func main() {
	app := &cli.App{
		Name:  "Biji",
		Usage: "Bitbucket & Jira CLI",
		Commands: []*cli.Command{
			{
				Name:    "repositories",
				Aliases: []string{"repo"},
				Usage:   "Show all repositories",
				Action: func(ctx *cli.Context) error {
					repositories := getRepositories()
					for i, repo := range repositories {
						fmt.Printf("%d. \033[1m%s\033[0m by %s\n", i+1, repo.FullName, repo.Owner.DisplayName)
					}
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("Welcome to Biji!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
