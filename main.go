package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := cli.App{
		Name:        "minicache",
		Usage:       "a mini cache server",
		Description: "A mini-cache server for golang simlar to redis but smaller",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listen",
				Aliases: []string{"l"},
				Usage:   "listen address",
				Value:   "0.0.0.0:8000",
			},
		},
		Action:  NewWithCtx,
		Version: Version,
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}

}
