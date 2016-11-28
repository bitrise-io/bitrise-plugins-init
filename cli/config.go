package cli

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-init/scanner"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v1"
)

var configCmd = cli.Command{
	Name:  "config",
	Usage: "Init bitrise config",
	Action: func(c *cli.Context) {
		if err := config(c); err != nil {
			log.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "empty",
			Usage: "creates empty bitrise config.",
		},
		cli.StringFlag{
			Name:  "platform",
			Usage: "creates bitrise config for a specified platform (options: ios,android,xamarin,fastlane)",
		},
		cli.StringFlag{
			Name:  "output",
			Usage: "bitrise config output file path",
			Value: "bitrise.yml",
		},
	},
}

func config(c *cli.Context) error {
	empty := c.Bool("empty")
	platform := c.String("platform")
	outputPth := c.String("output")

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory, error: %s", err)
	}

	if outputPth == "" {
		return fmt.Errorf("output not specified")
	}

	absOutputPth, err := pathutil.AbsPath(outputPth)
	if err != nil {
		return fmt.Errorf("failed to expand path (%s), error: %s", outputPth, err)
	}
	outputPth = absOutputPth

	bitriseConfig := bitriseModels.BitriseDataModel{}

	if empty {
		scanResult, err := scanner.ManualConfig()
		if err != nil {
			return fmt.Errorf("failed to create empty config, error: %s", err)
		}

		emptyConfigs, ok := scanResult.ConfigsMap["custom"]
		if !ok {
			return fmt.Errorf("no default empty configs found, error: %s", err)
		}

		emptyConfigStr, ok := emptyConfigs["custom-config"]
		if !ok {
			return fmt.Errorf("no default empty config found, error: %s", err)
		}

		var emptyConfig bitriseModels.BitriseDataModel
		if err := yaml.Unmarshal([]byte(emptyConfigStr), &emptyConfig); err != nil {
			return fmt.Errorf("no default empty config found, error: %s", err)
		}

		bitriseConfig = emptyConfig
	} else if platform != "" {
		scanResult, err := scanner.ManualConfig()
		if err != nil {
			return fmt.Errorf("failed to create empty config, error: %s", err)
		}

		platformOptions, ok := scanResult.OptionsMap[platform]
		if !ok {
			return fmt.Errorf("no available configs found for platform (%s), error: %s", platform, err)
		}

		configName, appEnvs, err := scanner.AskForOptions(platformOptions)
		if err != nil {
			return fmt.Errorf("failed to collect inputs, error: %s", err)
		}

		configMap := scanResult.ConfigsMap[platform]
		configStr := configMap[configName]

		var platformConfig bitriseModels.BitriseDataModel
		if err := yaml.Unmarshal([]byte(configStr), &platformConfig); err != nil {
			return fmt.Errorf("failed to unmarshal config, error: %s", err)
		}

		platformConfig.App.Environments = append(platformConfig.App.Environments, appEnvs...)

		bitriseConfig = platformConfig
	} else {
		scanResult, err := scanner.Config(currentDir)
		if err != nil {
			return err
		}

		if len(scanResult.OptionsMap) == 0 {
			return fmt.Errorf("No known platform type detected")
		}

		config, err := scanner.AskForConfig(scanResult)
		if err != nil {
			return err
		}

		bitriseConfig = config
	}

	if exist, err := pathutil.IsPathExists(outputPth); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("output path (%s) already exist", outputPth)
	}

	outputDir := filepath.Dir(outputPth)
	if exist, err := pathutil.IsDirExists(outputDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(outputDir, 0700); err != nil {
			return fmt.Errorf("failed to create output directory (%s), error: %s", outputDir, err)
		}
	}

	bytes, err := yaml.Marshal(bitriseConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal bitrise config, error: %s", err)
	}

	if err := fileutil.WriteBytesToFile(outputPth, bytes); err != nil {
		return fmt.Errorf("failed to write bitrise config, error: %s", err)
	}

	log.Infof("bitrise config generated at: %s", outputPth)

	return nil
}
