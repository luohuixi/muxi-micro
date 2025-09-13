package new

import (
	"github.com/muxi-Infra/muxi-micro/tool/new/create"
	"github.com/spf13/cobra"
)

func InitNewCobra() *cobra.Command {
	// new 子命令
	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "new 项目模板构建",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			if err := create.CreateAll(dir); err != nil {
				return err
			}
			return nil
		},
	}

	newCmd.Flags().String("dir", ".", "项目生成路径")

	return newCmd
}
