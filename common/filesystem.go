package common

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func FindFile(rootPath string, filter func(fs.DirEntry, string) (bool, string), maxDepth int) []string {
	resultDir := make([]string, 0)
	initDepth := strings.Count(rootPath, string(filepath.Separator))

	visit := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// Calculate the current depth by counting directory separators in the path
		depth := strings.Count(path, string(filepath.Separator)) - initDepth
		if depth > maxDepth {
			return filepath.SkipDir
		}

		if b, result := filter(d, path); b {
			resultDir = append(resultDir, result)
		}

		return nil
	}

	err := filepath.WalkDir(rootPath, visit)
	if err != nil {
		panic(err)
	}
	return resultDir
}
