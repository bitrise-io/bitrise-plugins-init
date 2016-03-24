package detectors

import (
	"os"
	"path/filepath"
	"strings"
)

func fileList(searchDir string) ([]string, error) {
	searchDir, err := filepath.Abs(searchDir)
	if err != nil {
		return []string{}, err
	}

	fileList := []string{}

	if err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)

		return nil
	}); err != nil {
		return []string{}, err
	}
	return fileList, nil
}

func filterFilesWithBasPaths(fileList []string, basePath ...string) []string {
	filteredFileList := []string{}

	for _, file := range fileList {
		base := filepath.Base(file)

		for _, desiredBasePath := range basePath {
			if strings.EqualFold(base, desiredBasePath) {
				filteredFileList = append(filteredFileList, file)
				break
			}
		}
	}

	return filteredFileList
}

func filterFilesWithExtensions(fileList []string, extension ...string) []string {
	filteredFileList := []string{}

	for _, file := range fileList {
		ext := filepath.Ext(file)

		for _, desiredExt := range extension {
			if ext == desiredExt {
				filteredFileList = append(filteredFileList, file)
				break
			}
		}
	}

	return filteredFileList
}

//--------------------------------------------------
// Sorting
//--------------------------------------------------

// ByComponents ..
type ByComponents []string

func (s ByComponents) Len() int {
	return len(s)
}
func (s ByComponents) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByComponents) Less(i, j int) bool {
	c1 := strings.Split(s[i], string(os.PathSeparator))
	c2 := strings.Split(s[j], string(os.PathSeparator))

	return len(c1) < len(c2)
}
