package create

import (
	"os"
	"path"
)

func CreateApi(output string) error {
	addr := path.Join(output, "example.api")
	content, _ := os.ReadFile("./gin/template/api.tpl")
	file, err := os.Create(addr)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(string(content))
	return nil
}
