package main

import (
	"biji/bitbucket"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

type Repository struct {
	FullName string `json:"full_name"`
	Owner    struct {
		DisplayName string `json:"display_name"`
	}
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

	bbUser := viper.GetString("bitbucket_username")
	bbPass := viper.GetString("bitbucket_password")
	client := bitbucket.NewClient(bbUser, bbPass)

	app := &cli.App{
		Name:  "Biji",
		Usage: "Bitbucket & Jira CLI",
		Commands: []*cli.Command{
			{
				Name:    "repositories",
				Aliases: []string{"repo"},
				Usage:   "Show all repositories",
				Action: func(ctx *cli.Context) error {
					repositories := client.GetRepositories()
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
			{
				Name:    "interactive",
				Usage:   "Open interactive mode",
				Aliases: []string{"ui"},
				Action: func(ctx *cli.Context) error {
					app := tview.NewApplication()
					box := tview.NewBox().SetBackgroundColor(tcell.ColorDefault).SetTitle("Box").SetBorder(true)
					flex := tview.NewFlex().AddItem(box, 0, 1, false)
					app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						if event.Rune() == 'q' {
							app.Stop()
						}
						return event
					})
					return app.SetRoot(flex, true).Run()
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
