package create

import (
	"os"
	"path"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

// 已存在就不覆盖
func CreateDocument(output, dir, addr string) ([]*parse.Api, error) {
	dir = path.Join(output, dir)
	dirs := []string{
		dir,
		path.Join(dir, "api"),
		path.Join(dir, "handler"),
		path.Join(dir, "router"),
		path.Join(dir, "router", "middleware"),
	}

	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			if err := os.MkdirAll(d, 0755); err != nil {
				return nil, err
			}
		}
	}

	if err := CopyAllApi(addr, path.Join(dir, "api")); err != nil {
		return nil, err
	}

	apis, err := parse.ParseAll(addr)
	if err != nil {
		return nil, err
	}

	return apis, nil
}
