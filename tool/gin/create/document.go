package create

import (
	"os"
	"path"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

// 已存在就不覆盖
func CreateDocument(output, addr string) ([]*parse.Api, error) {
	dirs := []string{
		output,
		path.Join(output, "handler"),
		path.Join(output, "router"),
	}

	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			if err := os.MkdirAll(d, 0755); err != nil {
				return nil, err
			}
		}
	}

	apis, err := parse.ParseAll(addr)
	if err != nil {
		return nil, err
	}

	return apis, nil
}
