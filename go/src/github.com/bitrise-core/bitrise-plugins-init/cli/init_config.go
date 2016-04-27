package cli

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-plugins-init/detectors"
	"github.com/bitrise-core/bitrise-plugins-init/models"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
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
	searchDir := "./"
	searchDir = "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-android"
	searchDir = "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-ios-cocoapods"
	searchDir = "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-xamarin-uitest"
	searchDir = "/Users/godrei/Develop/bitrise/sample-apps/fastlane-example"

	isCI := c.Bool("ci")
	isPrivate := c.Bool("private")
	log.Info("Configs:")
	log.Infof(" ci: %v", isCI)
	log.Infof(" private_repository: %v", isPrivate)
	fmt.Println()

	platformDetectors := []detectors.DetectorInterface{
		// new(detectors.Ios),
		// new(detectors.Android),
		// new(detectors.Xamarin),
		new(detectors.Fastlane),
	}

	optionsMap := map[string][]models.OptionModel{}
	configsMap := map[string]map[string]bitriseModels.BitriseDataModel{}

	// Run detectors
	log.Infof("Running platform detectors:")
	for _, detector := range platformDetectors {
		detectorName := detector.Name()
		fmt.Println()
		log.Infof("  Detector: %s", detectorName)

		detector.Configure(searchDir)

		detected, err := detector.DetectPlatform()
		if err != nil {
			log.Fatalf("Detector failed, error: %s", err)
		}

		if !detected {
			log.Info("  Platform not detected")
			continue
		}

		// Run analyzer
		log.Infof("  Running analyzer:")

		options, err := detector.Analyze()
		if err != nil {
			log.Fatalf("Analyzer failed, error: %s", err)
		}

		fmt.Println()
		log.Infof("  Analyze Result:")
		bytes, err := yaml.Marshal(options)
		if err != nil {
			log.Fatalf("Failed to marshal options, err: %s", err)
		}
		fmt.Printf("%v\n", string(bytes))

		optionsMap[detectorName] = options

		// Generate configs
		fmt.Println()
		log.Infof("  Generate configs:")
		configs := detector.Configs(isPrivate)
		for name, config := range configs {
			fmt.Printf("name: %s\n", name)
			bytes, err := yaml.Marshal(config)
			if err != nil {
				log.Fatalf("Failed to marshal options, err: %s", err)
			}
			fmt.Printf("%v\n", string(bytes))
		}

		configsMap[detectorName] = configs
	}

	// Write Env Options & Workflows to file
	if isCI {
		for detectorName, options := range optionsMap {
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
			optionsBytes, err := yaml.Marshal(options)
			if err != nil {
				log.Fatalf("Failed to marshal app envs, error: %s", err)
			}

			pth := path.Join(platformOutputDir, "app-envs.yml")
			if err := fileutil.WriteBytesToFile(pth, optionsBytes); err != nil {
				log.Fatalf("Failed to save app envs, err: %s", err)
			}
			log.Infof("app envs json saved to: %s", pth)

			// Bitrise Configs
			configMap := configsMap[detectorName]
			for configName, config := range configMap {
				configBytes, err := yaml.Marshal(config)
				if err != nil {
					log.Fatalf("Failed to marshal config, error: %#v", err)
				}

				pth = path.Join(platformOutputDir, configName)
				if err := fileutil.WriteBytesToFile(pth, configBytes); err != nil {
					log.Fatalf("Failed to save configs, err: %s", err)
				}
				log.Infof("bitrise.yml template saved to: %s", pth)
			}
		}

		return
	}

	// Collect app envs
	for detectorName, options := range optionsMap {
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

		fmt.Println()
		log.Infof("  Selected config: %s", configPth)
		log.Infof("  Selected envs:")
		aBytes, err := yaml.Marshal(appEnvs)
		if err != nil {
			log.Fatalf("Failed to marshal appEnvs, err: %s", err)
		}
		fmt.Printf("%v\n", string(aBytes))

		configMap := configsMap[detectorName]
		config := configMap[configPth]
		config.App.Environments = appEnvs

		log.Infof("  Selected config:")
		aBytes, err = yaml.Marshal(config)
		if err != nil {
			log.Fatalf("Failed to marshal config, err: %s", err)
		}
		fmt.Printf("%v\n", string(aBytes))

		// Write config to file
		configBytes, err := yaml.Marshal(config)
		if err != nil {
			log.Fatalf("Failed to marshal config, error: %#v", err)
		}

		pth := path.Join(platformOutputDir, configPth)
		if err := fileutil.WriteBytesToFile(pth, configBytes); err != nil {
			log.Fatalf("Failed to save configs, err: %s", err)
		}
		log.Infof("bitrise.yml template saved to: %s", pth)
	}
}
