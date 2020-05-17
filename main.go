package main

import (
	"fmt"
	"log"
	"os"

	"humanewolf.com/ed/systemapi/api"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "Serve the api.",
				Action: func(c *cli.Context) error {
					api.RunAPI()
					return nil
				},
			},
			{
				Name:    "prepare",
				Aliases: []string{"p"},
				Usage:   "Prepare the database.",
				Action: func(c *cli.Context) error {
					fmt.Printf("TODO, read: %s\n", c.Args().Get(0))
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
