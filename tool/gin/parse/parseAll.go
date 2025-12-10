package parse

import (
	"os"
	"path"
	"strings"
)

func ParseAll(sourceDir string) ([]*Api, error) {
	if _, err := os.Stat(sourceDir); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, err
	}

	var apis []*Api
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".api") {
			// 肯定存在
			api, err := ParesApi(path.Join(sourceDir, file.Name()))
			if err != nil {
				return nil, err
			}
			apis = append(apis, api)
		}
	}
	return apis, nil
}
