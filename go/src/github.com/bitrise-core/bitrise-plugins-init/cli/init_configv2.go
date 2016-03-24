package cli

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-plugins-init/detectors"
	"github.com/codegangsta/cli"
)

func initConfigV2(c *cli.Context) {
	if err := detectors.DetectAndroid(); err != nil {
		log.Fatalf("Android detector failed, error: %s", err)
	}

	if err := detectors.DetectIOS(); err != nil {
		log.Fatalf("iOS detector failed, error: %s", err)
	}
}
