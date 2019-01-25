package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vscreen/server-go/server"
)

func main() {
	var logLevel string

	app := cli.NewApp()
	app.Name = "vscreen"
	app.Usage = "realtime video remote"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "log-level, l",
			Value:       "info",
			Destination: &logLevel,
			Usage:       "set log level. [debug|info|error]",
		},
	}

	app.Action = func(c *cli.Context) error {
		// set log level
		switch logLevel {
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		default: // info
			log.SetLevel(log.InfoLevel)
		}

		s, err := server.New()
		if err != nil {
			return err
		}

		return s.ListenAndServe(":8080")
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
