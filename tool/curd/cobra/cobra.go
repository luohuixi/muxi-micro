package cobra

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

func curdCobra() {
	// curd 子命令
	var curdCmd = &cobra.Command{
		Use:   "curd",
		Short: "curd 自动生成工具",
		Run: func(cmd *cobra.Command, args []string) {
			pkg, err := cmd.Flags().GetString("package")
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {

			}
			err = CreateVar(pkg, dir)
			if err != nil {

			}
		},
	}

	curdCmd.Flags().String("package", "template", "生成文件的包名")
	curdCmd.Flags().String("dir", ".", "文件生成目录")

	//_ = curdCmd.Execute()
}

func CreateVar(pkg, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmplPath := filepath.Join("..", "template", "var.tpl")
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	t, err := template.New("var").Parse(string(tmplContent))
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dir, "var.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		PackageName string
	}{
		PackageName: pkg,
	}

	if err := t.Execute(file, data); err != nil {
		return err
	}

	return nil
}
