package utility

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// PackagesModel ...
type PackagesModel struct {
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func parsePackagesJSONContent(content string) (PackagesModel, error) {
	var packages PackagesModel
	if err := json.Unmarshal([]byte(content), &packages); err != nil {
		return PackagesModel{}, err
	}
	return packages, nil
}

// ParsePackagesJSON ...
func ParsePackagesJSON(packagesJSONPth string) (PackagesModel, error) {
	content, err := fileutil.ReadStringFromFile(packagesJSONPth)
	if err != nil {
		return PackagesModel{}, err
	}
	return parsePackagesJSONContent(content)
}

// CollectPackageJSONFiles ...
func CollectPackageJSONFiles(searchDir string) ([]string, error) {
	fileList, err := pathutil.ListPathInDirSortedByComponents(searchDir, false)
	if err != nil {
		return nil, err
	}

	filters := []pathutil.FilterFunc{
		pathutil.BaseFilter("package.json", true),
		pathutil.ComponentFilter("node_modules", false),
	}
	packageFileList, err := pathutil.FilterPaths(fileList, filters...)
	if err != nil {
		return nil, err
	}

	return packageFileList, nil
}

// RelPath ...
func RelPath(basePth, pth string) (string, error) {
	absBasePth, err := pathutil.AbsPath(basePth)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(absBasePth, "/private/var") {
		absBasePth = strings.TrimPrefix(absBasePth, "/private")
	}

	absPth, err := pathutil.AbsPath(pth)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(absPth, "/private/var") {
		absPth = strings.TrimPrefix(absPth, "/private")
	}

	return filepath.Rel(absBasePth, absPth)
}

// FileExists ...
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
