package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type Repository struct {
	FullName string `json:"full_name"`
	Owner    struct {
		DisplayName string `json:"display_name"`
	}
}

func getRepositories() []Repository {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.bitbucket.org/2.0/repositories", nil)
	if err != nil {
		panic(err)
	}

	bbUser := viper.GetString("bitbucket_username")
	bbPass := viper.GetString("bitbucket_password")
	if bbUser == "" || bbPass == "" {
		fmt.Println("Please login using `biji login` command.")
		os.Exit(1)
	}

	req.SetBasicAuth(bbUser, bbPass)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Invalid credentials")
		os.Exit(1)
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
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(home, ".bijirc")
	viper.SetConfigFile(path)
	viper.SetConfigType("json")
	viper.ReadInConfig()

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
			{
				Name:  "login",
				Usage: "Login to bitbucket",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u"},
						Usage:   "Bitbucket username",
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
						Usage:   "Bitbucket app password",
					},
				},
				Action: func(ctx *cli.Context) error {
					viper.Set("bitbucket_username", ctx.String("username"))
					viper.Set("bitbucket_password", ctx.String("password"))

					if err := viper.WriteConfigAs(path); err != nil {
						fmt.Println(err)
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
