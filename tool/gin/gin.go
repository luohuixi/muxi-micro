package gin

import (
	"github.com/muxi-Infra/muxi-micro/tool/gin/create"
	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
	"github.com/spf13/cobra"
)

func InitGinCobra() *cobra.Command {
	// gin 子命令
	var ginCmd = &cobra.Command{
		Use:   "gin",
		Short: "gin 生成gin框架",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("api")
			output, _ := cmd.Flags().GetString("output")
			api, err := parse.ParesApi(addr)
			if err != nil {
				return err
			}
			if err := create.CreateService(output, api); err != nil {
				return err
			}
			return nil
		},
	}

	var ginAllCmd = &cobra.Command{
		Use:   "all",
		Short: "根据所有的api直接生成整个gin框架",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, _ := cmd.Flags().GetString("api")
			output, _ := cmd.Flags().GetString("output")
			dir, _ := cmd.Flags().GetString("dir")
			apis, err := create.CreateDocument(output, dir, addr)
			if err != nil {
				return err
			}
			if err := create.CreateAllService(output, dir, apis); err != nil {
				return err
			}
			return nil
		},
	}

	var ginApiCmd = &cobra.Command{
		Use:   "api",
		Short: "给一个api文件的示例, 就不用记住格式了",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			if err := create.CreateApi(output); err != nil {
				return err
			}
			return nil
		},
	}

	ginCmd.Flags().String("api", ".", "api文件的位置")
	ginCmd.Flags().String("output", ".", "输出路径")

	ginAllCmd.Flags().String("api", ".", "存放所有api文件的目录")
	ginAllCmd.Flags().String("output", ".", "输出路径")
	ginAllCmd.Flags().String("dir", "gin", "文件夹名称")

	ginApiCmd.Flags().String("output", ".", "输出路径")

	ginCmd.AddCommand(ginAllCmd)
	ginCmd.AddCommand(ginApiCmd)

	return ginCmd
}
