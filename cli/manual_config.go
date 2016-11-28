package cli

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-init/scanner"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v1"
)

var manualConfigCmd = cli.Command{
	Name:  "manual-config",
	Usage: "Creates bitrise template config and secrets",
	Action: func(c *cli.Context) {
		if err := manualConfig(c); err != nil {
			log.Fatal(err)
		}
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "empty",
			Usage: "creates empty bitrise config and secrets.",
		},
		cli.StringFlag{
			Name:  "platform",
			Usage: "creates bitrise config for a specified platform (options: ios,android,xamarin,fastlane) and secrets",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "bitrise config file path",
			Value: "bitrise.yml",
		},
		cli.StringFlag{
			Name:  "secrets",
			Usage: "bitrise secrets file path",
			Value: ".bitrise.secrets.yml",
		},
	},
}

func manualConfig(c *cli.Context) error {
	// validate inputs
	configPth := c.String("config")
	secretsPth := c.String("secrets")
	empty := c.Bool("empty")
	platform := c.String("platform")

	if configPth == "" {
		return fmt.Errorf("config path not specified")
	}

	if exist, err := pathutil.IsPathExists(configPth); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("config path (%s) already exist", configPth)
	}

	if secretsPth == "" {
		return fmt.Errorf("secrets path not specified")
	}

	if exist, err := pathutil.IsPathExists(secretsPth); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("secrets path (%s) already exist", secretsPth)
	}

	if !c.IsSet("empty") && !c.IsSet("platform") {
		return fmt.Errorf("Use --empty or --platform (ios/android/xamarin/fastlane) flag to specify config type")
	}

	// create template config
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
	}

	// write outputs
	configDir := filepath.Dir(configPth)
	if exist, err := pathutil.IsDirExists(configDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("failed to create config directory (%s), error: %s", configDir, err)
		}
	}

	secretsDir := filepath.Dir(secretsPth)
	if exist, err := pathutil.IsDirExists(secretsDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(secretsDir, 0700); err != nil {
			return fmt.Errorf("failed to create secrets directory (%s), error: %s", secretsDir, err)
		}
	}

	configBytes, err := yaml.Marshal(bitriseConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal bitrise config, error: %s", err)
	}

	if err := fileutil.WriteBytesToFile(configPth, configBytes); err != nil {
		return fmt.Errorf("failed to write bitrise config, error: %s", err)
	}

	log.Infof("bitrise config generated at: %s", configPth)

	secrets := envmanModels.EnvsYMLModel{}

	secretsBytes, err := yaml.Marshal(secrets)
	if err != nil {
		return fmt.Errorf("failed to marshal bitrise secrets, error: %s", err)
	}

	if err := fileutil.WriteBytesToFile(secretsPth, secretsBytes); err != nil {
		return fmt.Errorf("failed to write bitrise secrets, error: %s", err)
	}

	log.Infof("bitrise secrets generated at: %s", secretsPth)

	return nil
}
