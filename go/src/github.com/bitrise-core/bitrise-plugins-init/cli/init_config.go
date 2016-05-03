package cli

import (
	"fmt"
	"os"
	"path"

	"encoding/json"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-plugins-init/detectors"
	"github.com/bitrise-core/bitrise-plugins-init/models"
	"github.com/bitrise-core/bitrise-plugins-init/version"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/codegangsta/cli"
)

const outputDir = "/Users/godrei/Develop/bitrise/plugins/bitrise-plugins-init/go/src/github.com/bitrise-core/bitrise-plugins-init/_tmp"

func askForValue(option models.OptionModel) (string, string, error) {
	optionValues := option.GetValues()

	selectedValue := ""
	if len(optionValues) == 1 {
		selectedValue = optionValues[0]
	} else {
		question := fmt.Sprintf("Select: %s (%s)", option.Title, option.Key)
		answer, err := goinp.SelectFromStrings(question, optionValues)
		if err != nil {
			return "", "", err
		}

		selectedValue = answer
	}

	return option.EnvKey, selectedValue, nil
}

func initConfig(c *cli.Context) {
	//
	// Config
	isCI := c.GlobalBool("ci")
	isPrivate := c.Bool("private")
	searchDir := c.String("dir")
	if searchDir == "" {
		searchDir = "./"
		// searchDir = "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-ios-cocoapods"
		searchDir = "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-android"
		// searchDir = "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-xamarin-uitest"
		// searchDir = "/Users/godrei/Develop/bitrise/sample-apps/fastlane-example"
	}

	fmt.Println()
	log.Info(colorstring.Greenf("Running scanner v%s", version.VERSION))
	fmt.Println()

	if isCI {
		log.Info(colorstring.Yellow("plugin runs in CI mode"))
	}
	if isPrivate {
		log.Info(colorstring.Yellow("scanning private repository"))
	}
	log.Info(colorstring.Yellowf("scan dir: %s", searchDir))
	fmt.Println()

	//
	// Scan
	platformDetectors := []detectors.DetectorInterface{
		new(detectors.Android),
		new(detectors.Xamarin),
		new(detectors.Ios),
		new(detectors.Fastlane),
	}
	optionsMap := map[string][]models.OptionModel{}
	configsMap := map[string]map[string]bitriseModels.BitriseDataModel{}

	log.Infof(colorstring.Blue("Running scanners:"))
	for _, detector := range platformDetectors {
		detectorName := detector.Name()
		log.Infof("  Scanner: %s", colorstring.Blue(detectorName))

		detector.Configure(searchDir)
		detected, err := detector.DetectPlatform()
		if err != nil {
			log.Fatalf("Scanner failed, error: %s", err)
		}

		if !detected {
			log.Info("  Platform not detected")
			fmt.Println()
			continue
		}

		log.Info("  Platform detected")
		log.Info("  +----------------------------------------+")
		log.Info("  |                                        |")

		options, err := detector.Analyze()
		if err != nil {
			log.Fatalf("Analyzer failed, error: %s", err)
		}

		log.Debug()
		log.Debug("Analyze result:")
		bytes, err := yaml.Marshal(options)
		if err != nil {
			log.Fatalf("Failed to marshal options, err: %s", err)
		}
		log.Debugf("\n%v", string(bytes))

		optionsMap[detectorName] = options

		// Generate configs
		log.Debug()
		log.Debug("Generated configs:")
		configs := detector.Configs(isPrivate)
		for name, config := range configs {
			log.Debugf("  name: %s", name)

			bytes, err := yaml.Marshal(config)
			if err != nil {
				log.Fatalf("Failed to marshal options, err: %s", err)
			}
			log.Debugf("\n%v", string(bytes))
		}

		configsMap[detectorName] = configs

		log.Info("  |                                        |")
		log.Info("  +----------------------------------------+")
		fmt.Println()
	}

	//
	// Write output to files
	if isCI {
		log.Infof(colorstring.Blue("Saving outputs:"))

		for detectorName, options := range optionsMap {
			log.Infof("  Scanner: %s", colorstring.Blue(detectorName))

			// Init
			platformOutputDir := path.Join(outputDir, detectorName)
			if exist, err := pathutil.IsDirExists(platformOutputDir); err != nil {
				log.Fatalf("Failed to check if path (%s) exis, err: %s", platformOutputDir, err)
			} else if exist {
				if err := os.RemoveAll(platformOutputDir); err != nil {
					log.Fatalf("Failed to cleanup (%s), err: %s", platformOutputDir, err)
				}
			}

			if err := os.MkdirAll(platformOutputDir, 0700); err != nil {
				log.Fatalf("Failed to create (%s), err: %s", platformOutputDir, err)
			}

			// App Envs Options
			optionsBytes, err := json.Marshal(options)
			if err != nil {
				log.Fatalf("Failed to marshal app envs, error: %s", err)
			}

			pth := path.Join(platformOutputDir, "app-envs.json")
			if err := fileutil.WriteBytesToFile(pth, optionsBytes); err != nil {
				log.Fatalf("Failed to save app envs, err: %s", err)
			}
			log.Infof("  app env options: %s", colorstring.Blue(pth))

			// Bitrise Configs
			configMap := configsMap[detectorName]
			for configName, config := range configMap {
				configBytes, err := json.Marshal(config)
				if err != nil {
					log.Fatalf("Failed to marshal config, error: %#v", err)
				}

				pth = path.Join(platformOutputDir, configName)
				if err := fileutil.WriteBytesToFile(pth, configBytes); err != nil {
					log.Fatalf("Failed to save configs, err: %s", err)
				}
				log.Infof("  bitrise.json template: %s", colorstring.Blue(pth))
			}

			fmt.Println()
		}

		return
	}

	//
	// Select options
	log.Infof(colorstring.Blue("Collecting inputs:"))

	for detectorName, options := range optionsMap {
		log.Infof("  Scanner: %s", colorstring.Blue(detectorName))

		// Init
		platformOutputDir := path.Join(outputDir, detectorName)
		if exist, err := pathutil.IsDirExists(platformOutputDir); err != nil {
			log.Fatalf("Failed to check if path (%s) exis, err: %s", platformOutputDir, err)
		} else if exist {
			if err := os.RemoveAll(platformOutputDir); err != nil {
				log.Fatalf("Failed to cleanup (%s), err: %s", platformOutputDir, err)
			}
		}

		if err := os.MkdirAll(platformOutputDir, 0700); err != nil {
			log.Fatalf("Failed to create (%s), err: %s", platformOutputDir, err)
		}

		// Collect inputs
		configPth := ""
		appEnvs := []envmanModels.EnvironmentItemModel{}

		var walkWidth func(options []models.OptionModel)

		walkDepth := func(option models.OptionModel) {
			optionEnvKey, selectedValue, err := askForValue(option)
			if err != nil {
				log.Fatalf("Failed to ask for vale of key (%s), error: %s", option.Key, err)
			}

			if optionEnvKey == "" {
				configPth = selectedValue
			} else {
				appEnvs = append(appEnvs, envmanModels.EnvironmentItemModel{
					optionEnvKey: selectedValue,
				})
			}

			nestedOptions := option.ValueMap[selectedValue]
			if len(nestedOptions) == 0 {
				return
			}

			walkWidth(nestedOptions)
		}

		walkWidth = func(options []models.OptionModel) {
			for _, option := range options {
				walkDepth(option)
			}
		}

		walkWidth(options)

		log.Debug()
		log.Debug("Selected app envs:")
		aBytes, err := yaml.Marshal(appEnvs)
		if err != nil {
			log.Fatalf("Failed to marshal appEnvs, err: %s", err)
		}
		log.Debugf("\n%v", string(aBytes))

		configMap := configsMap[detectorName]
		config := configMap[configPth]
		config.App.Environments = appEnvs

		log.Debug()
		log.Debug("Config:")
		log.Debugf("  name: %s", configPth)
		aBytes, err = yaml.Marshal(config)
		if err != nil {
			log.Fatalf("Failed to marshal config, err: %s", err)
		}
		log.Debugf("\n%v", string(aBytes))

		// Write config to file
		configBytes, err := json.Marshal(config)
		if err != nil {
			log.Fatalf("Failed to marshal config, error: %#v", err)
		}

		pth := path.Join(platformOutputDir, configPth)
		if err := fileutil.WriteBytesToFile(pth, configBytes); err != nil {
			log.Fatalf("Failed to save configs, err: %s", err)
		}
		log.Infof("  bitrise.json template: %s", colorstring.Blue(pth))
		fmt.Println()
	}
}
