package main

import (
	"fmt"
	"os"

	"github.com/h3th-IV/mackerel/internal/command"
	"github.com/urfave/cli/v2"
)

func main() {
	fmt.Println("hello world")
	app := &cli.App{
		Name:  "mackerel",
		Usage: "A phishing simulation tool",
		Commands: []*cli.Command{
			command.StartCommand(),
		},
		Version: "v0.1.5",
		Authors: []*cli.Author{
			{
				Name:  "h3th-IV",
				Email: "h3th-IV@proton.mail",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error running programs:", err.Error())
		os.Exit(1)
	}
}
