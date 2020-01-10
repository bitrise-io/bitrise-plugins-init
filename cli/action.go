package cli

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-init/scanner"
	"github.com/bitrise-io/bitrise-init/scanners"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/urfave/cli"
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

		customConfigs, ok := scanResult.ScannerToBitriseConfigMap[scanners.CustomProjectType]
		if !ok {
			return fmt.Errorf("no CustomProjectType found found, error: %s", err)
		}

		customConfigStr, ok := customConfigs[scanners.CustomConfigName]
		if !ok {
			return fmt.Errorf("no CustomConfig found, error: %s", err)
		}

		var customConfig bitriseModels.BitriseDataModel
		if err := yaml.Unmarshal([]byte(customConfigStr), &customConfig); err != nil {
			return fmt.Errorf("no default empty config found, error: %s", err)
		}

		bitriseConfig = customConfig
	} else {
		// run scanner
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory, error: %s", err)
		}

		scanResult := scanner.Config(currentDir)

		if len(scanResult.ScannerToOptionRoot) == 0 {
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

	secrets := envmanModels.EnvsSerializeModel{}
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

func gitignore(pattern, gitignorePath string) error {
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open .gitignore file at %s: %s", gitignorePath, err)
	}

	if _, err := f.WriteString("\n" + pattern); err != nil {
		return fmt.Errorf("write pattern to .gitignore at %s: %s", gitignorePath, err)
	}

	return nil
}
