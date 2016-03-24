package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-plugins-init/detectors"
	"github.com/codegangsta/cli"
)

func initConfigV2(c *cli.Context) {
	fmt.Println()
	log.Infof("----- Running android detector -----")
	if err := detectors.DetectAndroid(); err != nil {
		log.Fatalf("Android detector failed, error: %s", err)
	}
	log.Infof("------------------------------------")

	fmt.Println()
	log.Infof("----- Running iOS detector -----")
	if err := detectors.DetectIOS(); err != nil {
		log.Fatalf("iOS detector failed, error: %s", err)
	}
	log.Infof("--------------------------------")

	fmt.Println()
	log.Infof("----- Running xamarin detector -----")
	if err := detectors.DetectXamarin(); err != nil {
		log.Fatalf("Xamarin detector failed, error: %s", err)
	}
	log.Infof("------------------------------------")
}
