package detectors

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/fileutil"
)

const (
	solutionExtension          = ".sln"
	solutionConfigurationStart = "GlobalSection(SolutionConfigurationPlatforms) = preSolution"
	solutionConfigurationEnd   = "EndGlobalSection"

	includeMonoTouchAPIPattern   = `Include="monotouch"`
	includeXamarinIosAPIPattern  = `Include="Xamarin.iOS"`
	includeMonoAndroidAPIPattern = `Include="Mono.Android`

	monoTouchAPI   = "monotouch"
	xamarinIosAPI  = "Xamarin.iOS"
	monoAndroidAPI = "Mono.Android"
)

//--------------------------------------------------
// Utility
//--------------------------------------------------

func filterSolutionFiles(fileList []string) []string {
	return filterFilesWithExtensions(fileList, solutionExtension)
}

func getSolutionConfigs(solutionFile string) ([]string, error) {
	content, err := fileutil.ReadStringFromFile(solutionFile)
	if err != nil {
		return []string{}, err
	}

	configs := []string{}
	isNextLineScheme := false

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, solutionConfigurationStart) {
			isNextLineScheme = true
			continue
		}

		if strings.Contains(line, solutionConfigurationEnd) {
			isNextLineScheme = false
			continue
		}

		if isNextLineScheme {
			split := strings.Split(line, "=")
			if len(split) == 2 {
				config := strings.TrimSpace(split[1])
				configs = append(configs, config)
			} else {
				log.Warnf("Invalid config: %s", line)
			}
		}
	}

	return configs, nil
}

func getProjectPlatformAPI(projectFile string) (string, error) {
	content, err := fileutil.ReadStringFromFile(projectFile)
	if err != nil {
		return "", err
	}

	if saseInsensitiveContains(content, includeMonoAndroidAPIPattern) {
		return monoAndroidAPI, nil
	} else if saseInsensitiveContains(content, includeMonoTouchAPIPattern) {
		return monoTouchAPI, nil
	} else if saseInsensitiveContains(content, includeXamarinIosAPIPattern) {
		return xamarinIosAPI, nil
	}

	return "", nil
}

func getProjects(solutionFile string) ([]string, error) {
	content, err := fileutil.ReadStringFromFile(solutionFile)
	if err != nil {
		return []string{}, err
	}

	projectDir := filepath.Dir(solutionFile)
	projectExp := regexp.MustCompile(`Project\(\"[^\"]*\"\)\s*=\s*\"[^\"]*\",\s*\"([^\"]*.csproj)\"`)

	projects := []string{}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		matches := projectExp.FindStringSubmatch(line)
		if len(matches) > 1 {
			project := strings.Replace(matches[1], "\\", "/", -1)
			projectPath := path.Join(projectDir, project)
			projects = append(projects, projectPath)
		}
	}

	return projects, nil
}

//--------------------------------------------------
// Main
//--------------------------------------------------

// DetectXamarin ...
func DetectXamarin() error {
	searchDir := "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-xamarin-uitest"
	// searchDir = "./"

	fileList, err := fileList(searchDir)
	if err != nil {
		return fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}

	// Search for solution file
	solutionFiles := filterSolutionFiles(fileList)

	log.Debugf("%s files:", solutionExtension)
	for _, solutionFile := range solutionFiles {
		log.Debugf("  %s", solutionFile)
	}
	log.Debugln("")

	if len(solutionFiles) == 0 {
		log.Infof("NO xamarin solution detected")
		return nil
	}

	log.Infof("xamarin solution detected")
	fmt.Println()

	// Check for solution configs
	validSolutionMap := map[string][]string{}
	for _, solutionFile := range solutionFiles {
		configs, err := getSolutionConfigs(solutionFile)
		if err != nil {
			return err
		}

		if len(configs) > 0 {
			validSolutionMap[solutionFile] = configs
		} else {
			log.Warnf("No config found for %s", solutionFile)
		}
	}

	// Check for solution projects
	configMap := map[string]map[string]interface{}{}

	for solutionFile, configs := range validSolutionMap {
		projects, err := getProjects(solutionFile)
		if err != nil {
			return err
		}

		solutionConfigMap := map[string]interface{}{
			"configs": configs,
		}

		// Inspect projects
		apis := []string{}
		for _, project := range projects {
			log.Infof("Inspecting project file: %s", project)

			api, err := getProjectPlatformAPI(project)
			if err != nil {
				return err
			}

			if api == "" {
				continue
			}

			apis = append(apis, api)
		}

		solutionConfigMap["apis"] = apis

		configMap[solutionFile] = solutionConfigMap
	}
	fmt.Println()

	// Log configurations
	for solutionFile, solutionConfigs := range configMap {
		log.Infof("Configurations for %s", solutionFile)
		for key, value := range solutionConfigs {
			log.Infof("  %s: %v", key, value)
		}
		fmt.Println()
	}

	return nil
}
