package cli

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-init/scanner"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/depman/pathutil"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v1"
)

func action(c *cli.Context) error {
	minimal := c.Bool("minimal")

	configPth := "./bitrise.yml"
	if exist, err := pathutil.IsPathExists(configPth); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("config path (%s) already exist", configPth)
	}

	secretsPth := "./.bitrise.secrets.yml"
	if exist, err := pathutil.IsPathExists(secretsPth); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("secrets path (%s) already exist", secretsPth)
	}

	// generate config
	bitriseConfig := bitriseModels.BitriseDataModel{}
	if minimal {
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
	} else {
		// run scanner
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory, error: %s", err)
		}

		scanResult, err := scanner.Config(currentDir)
		if err != nil {
			return err
		}

		if len(scanResult.OptionsMap) == 0 {
			return fmt.Errorf("no known platform type detected")
		}

		config, err := scanner.AskForConfig(scanResult)
		if err != nil {
			return err
		}

		bitriseConfig = config
	}

	// write outputs
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
