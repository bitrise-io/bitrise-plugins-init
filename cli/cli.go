package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/bitrise-io/bitrise-plugins-init/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Run ...
func Run() {
	cli.VersionPrinter = func(c *cli.Context) { fmt.Println(c.App.Version) }

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})

	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "Init bitrise config"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Action = func(c *cli.Context) error {
		if err := action(c); err != nil {
			log.Fatal(err)
		}

		return nil
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "minimal",
			Usage: "create empty bitrise config and secrets",
		},
		cli.BoolFlag{
			Name:  "private",
			Usage: "is a private repository",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
