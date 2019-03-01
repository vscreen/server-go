package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	vplayer "github.com/vscreen/server-go/player"
	"github.com/vscreen/server-go/server"
)

func main() {
	var logLevel string
	var player string
	players := strings.Join(vplayer.Players, "|")

	app := cli.NewApp()
	app.Name = "vscreen"
	app.Usage = "realtime video remote"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "log-level, l",
			Value:       "info",
			Destination: &logLevel,
			Usage:       "set log level [debug|info|error]",
		},
		cli.StringFlag{
			Name:        "player, p",
			Value:       "vlc",
			Destination: &player,
			Usage:       fmt.Sprintf("set player interface [%s]", players),
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

		p, err := vplayer.New(player)
		if err != nil {
			return err
		}
		defer p.Close()

		s, err := server.New(p)
		if err != nil {
			return err
		}

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		publish(ctx)

		return s.ListenAndServe(fmt.Sprintf(":%d", port))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
