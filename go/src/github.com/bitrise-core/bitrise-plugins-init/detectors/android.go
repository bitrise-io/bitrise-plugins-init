package detectors

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
)

const (
	buildGradleBasePath = "build.gradle"
	gradlewBasePath     = "gradlew"
)

//--------------------------------------------------
// Utility
//--------------------------------------------------

func filterGradleFiles(fileList []string) []string {
	filteredFileList := filterFilesWithBasPaths(fileList, buildGradleBasePath)
	sort.Sort(ByComponents(filteredFileList))
	return filteredFileList
}

func filterGradlewFiles(fileList []string) []string {
	filteredFileList := filterFilesWithBasPaths(fileList, gradlewBasePath)
	sort.Sort(ByComponents(filteredFileList))
	return filteredFileList
}

func inspectGradleFile(gradleFile string, gradleBin string) ([]string, error) {
	out, err := cmdex.RunCommandAndReturnCombinedStdoutAndStderr(gradleBin, "tasks", "--build-file", gradleFile)
	if err != nil {
		return []string{}, fmt.Errorf("output: %s, error: %s", out, err)
	}

	lines := strings.Split(out, "\n")
	isBuildTaskSection := false
	buildTasksExp := regexp.MustCompile(`^Build tasks`)
	configurationExp := regexp.MustCompile(`^(assemble\S+)(\s*-\s*.*)*`)

	configurations := []string{}
	for _, line := range lines {
		if !isBuildTaskSection && buildTasksExp.FindString(line) != "" {
			isBuildTaskSection = true
			continue
		} else if line == "" {
			isBuildTaskSection = false
			continue
		}

		if !isBuildTaskSection {
			continue
		}

		match := configurationExp.FindStringSubmatch(line)
		if len(match) > 1 {
			configurations = append(configurations, match[1])
		}
	}

	return configurations, nil
}

//--------------------------------------------------
// Main
//--------------------------------------------------

// DetectAndroid ...
func DetectAndroid() error {
	searchDir := "/Users/godrei/Develop/bitrise/sample-apps/sample-apps-android"
	// searchDir = "./"

	fileList, err := fileList(searchDir)
	if err != nil {
		return fmt.Errorf("failed to search for files in (%s), error: %s", searchDir, err)
	}

	// Search for gradle file
	gradleFiles := filterGradleFiles(fileList)

	log.Debugf("%s files:", buildGradleBasePath)
	for _, gradleFile := range gradleFiles {
		log.Debugf("  %s", gradleFile)
	}
	log.Debugln("")

	if len(gradleFiles) == 0 {
		log.Infof("NO android gradle project detected")
		return nil
	}

	log.Infof("android gradle project detected")
	fmt.Println()

	// Search for gradlew file
	gradlewFiles := filterGradlewFiles(fileList)

	log.Debugf("%s files:", gradlewBasePath)
	for _, gradlewFile := range gradlewFiles {
		log.Debugf("  %s", gradlewFile)
	}
	log.Debugln("")

	if len(gradlewFiles) == 0 {
		log.Infof("No gradlew file detected")
	}

	// Most root gradlew file
	rootGradlewPath := ""
	if len(gradlewFiles) > 0 {
		rootGradlewPath = gradlewFiles[0]
	}

	// Gradle or Gradlew
	gradleBin := "gradle"
	if rootGradlewPath != "" {
		log.Infof("Gradle wrapper (gradlew) found - using it: %s", rootGradlewPath)

		err := os.Chmod(rootGradlewPath, 0770)
		if err != nil {
			return fmt.Errorf("failed to add executable permission on gradlew file (%s), error: %s", rootGradlewPath, err)
		}

		gradleBin = rootGradlewPath
	} else {
		log.Info("No gradle wrapper (gradlew) found - using gradle directly")
	}
	fmt.Println()

	// Inspect Gradle files
	configMap := map[string][]string{}

	for _, gradleFile := range gradleFiles {
		log.Infof("Inspecting gradle file: %s", gradleFile)
		log.Infof(" $ %s tasks --build-file %s", gradleBin, gradleFile)

		configs, err := inspectGradleFile(gradleFile, gradleBin)
		if err != nil {
			return fmt.Errorf("failed to inspect gradle files, error: %s", err)
		}

		configMap[gradleFile] = configs
	}
	fmt.Println()

	// Log configurations
	for path, configs := range configMap {
		log.Infof("Configurations for %s", path)
		log.Info("Assemble configurations (except 'assemble'):")
		for _, config := range configs {
			log.Infof("  %s", config)
		}
		fmt.Println()
	}

	return nil
}
