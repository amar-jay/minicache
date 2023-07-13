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
		Usage:       "",
		Description: "",
		Flags:       []cli.Flag{},
		Action:      NewWithCtx,
		Version:     Version,
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n", c.App.Version)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}

}
