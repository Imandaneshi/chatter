package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name:     "Chatter",
		Usage:    "Chatter golang web server",
		Compiled: time.Now(),
		Version:  "0.1",
		Authors: []*cli.Author{
			{
				Name:  "Iman Daneshi",
				Email: "emandaneshikohan@gmail.com",
			},
		},

		Commands: []*cli.Command{
			Server(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("failed starting the socket server")
	}
}
