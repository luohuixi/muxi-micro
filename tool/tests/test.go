package tests

import (
	"github.com/muxi-Infra/muxi-micro/tool/tests/create"
	"github.com/muxi-Infra/muxi-micro/tool/tests/parse"
	"github.com/spf13/cobra"
)

func InitTestCobra() *cobra.Command {
	// test 子命令
	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "test 生成测试文件",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, _ := cmd.Flags().GetString("file")
			function, _ := cmd.Flags().GetBool("func")
			receiver, _ := cmd.Flags().GetBool("receive")
			funcStruct, RecesStruct, err := parse.ParseFunc(filePath)
			if err != nil {
				return err
			}
			packName, err := parse.ParsePackage(filePath)
			if err != nil {
				return err
			}
			if function {
				if err := create.CreateFunc(filePath, funcStruct, packName); err != nil {
					return err
				}
			}
			if receiver {
				if err := create.CreateRece(filePath, RecesStruct, packName); err != nil {
					return err
				}
			}
			return nil
		},
	}

	testCmd.Flags().String("file", "", "想要生成测试文件的路径")
	testCmd.Flags().Bool("func", false, "生成函数的测试")
	testCmd.Flags().Bool("receive", false, "生成方法的测试")
	testCmd.MarkFlagRequired("file")

	return testCmd
}
