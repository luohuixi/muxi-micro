package cobra

import (
	"fmt"
	"github.com/muxi-Infra/muxi-micro/tool/curd/create"
	"github.com/muxi-Infra/muxi-micro/tool/curd/parse"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func InitCurdCobra(parentCobra *cobra.Command) {
	// curd 子命令
	var curdCmd = &cobra.Command{
		Use:   "curd",
		Short: "curd 自动生成工具",
		Long: "在你想要生成 curd 文件的地方创建 model.go，内含可通过 gorm 自动迁移的结构体，" +
			"gorm 标签中 'primaryKey' 将被视为主键，主键和gorm标签中带有 'unique' 'index' 会自动生成查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmd.Flags().GetString("package")
			if err != nil {
				return err
			}
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				return err
			}

			modelPath := filepath.Join(dir, "model.go")
			if _, err := os.Stat(modelPath); err != nil {
				return err
			}
			table, index, err := parse.ParseStruct(modelPath)
			if err != nil {
				return err
			}

			err = create.CreateExample_gen(pkg, dir, table, index)
			if err != nil {
				return err
			}
			err = create.CreateExample(pkg, dir, table)
			if err != nil {
				return err
			}
			err = create.CreateVar(pkg, dir)
			if err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}

	curdCmd.Flags().String("package", "template", "生成文件的包名")
	curdCmd.Flags().String("dir", ".", "model文件以及文件生成目录")

	parentCobra.AddCommand(curdCmd)
}
