package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

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

	if err := gitignore(".bitrise.secrets.yml", "./.gitignore"); err != nil {
		log.Warnf("Could not add .bitrise.secrets.yml to .gitignore: %s", err)
		log.Warnf("Please be advised, that for security considerations, it is not recommended to upload .bitrise.secrets.yml to version control")
	}

	return nil
}

func gitignore(pattern, gitignorePath string) error {
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open .gitignore file at %s: %s", gitignorePath, err)
	}
	defer f.Close()

	contents, err := ioutil.ReadFile(gitignorePath)
	matched, err := regexp.MatchString(fmt.Sprintf("^%s$", pattern), string(contents))
	if err != nil {
		return fmt.Errorf("matching .gitignore file contents at %s against %s: %s", gitignorePath, pattern, err)
	}
	if matched {
		return nil
	}

	if len(contents) > 0 && !strings.HasSuffix(string(contents), fmt.Sprintln("")) {
		pattern = "\n" + pattern
	}

	if _, err := f.WriteString(pattern); err != nil {
		return fmt.Errorf("write pattern to .gitignore at %s: %s", gitignorePath, err)
	}

	return nil
}
