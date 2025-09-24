package curd

import (
	"os"
	"path/filepath"

	"github.com/muxi-Infra/muxi-micro/tool/curd/create"
	"github.com/muxi-Infra/muxi-micro/tool/curd/parse"
	"github.com/spf13/cobra"
)

func InitCurdCobra() *cobra.Command {
	// curd 子命令
	var curdCmd = &cobra.Command{
		Use:   "curd",
		Short: "curd 自动生成工具",
		Long: "在你想要生成 curd 文件的地方创建 model.go，内含可通过 gorm 自动迁移的结构体，" +
			"gorm 标签中 'primaryKey' 将被视为主键，主键和gorm标签中带有 'unique' 'index' 会自动生成查询",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, _ := cmd.Flags().GetString("package")
			dir, _ := cmd.Flags().GetString("dir")
			cache, _ := cmd.Flags().GetBool("cache")
			cover, _ := cmd.Flags().GetBool("cover")
			transcation, _ := cmd.Flags().GetBool("transcation")

			modelPath := filepath.Join(dir, "model.go")
			if _, err := os.Stat(modelPath); err != nil {
				return err
			}
			structs, err := parse.ParseStructs(modelPath)
			if err != nil {
				return err
			}

			for _, v := range structs {
				var index []parse.FieldInfo
				for _, vv := range v.Index {
					index = append(index, vv)
				}
				index = append(index, v.Primary[0])
				err = create.CreateExample_gen(pkg, dir, v.Name, index, cache)
				if err != nil {
					return err
				}
				err = create.CreateExample(pkg, dir, v.Name, cache, cover)
				if err != nil {
					return err
				}
				err = create.CreateVar(pkg, v.Name, dir, cache, cover)
				if err != nil {
					return err
				}
				err = create.CreateTranscation(pkg, dir, transcation)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	curdCmd.Flags().String("package", "repository", "生成文件的包名")
	curdCmd.Flags().String("dir", ".", "model文件以及文件生成目录")
	curdCmd.Flags().Bool("cache", false, "是否开启缓存")
	curdCmd.Flags().Bool("cover", false, "是否覆盖除 _gen.go 外的另外两个文件")
	curdCmd.Flags().Bool("transcation", false, "是否生成事务代码")

	return curdCmd
}
