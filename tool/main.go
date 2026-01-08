package main

import (
	"github.com/muxi-Infra/muxi-micro/tool/curd"
	"github.com/muxi-Infra/muxi-micro/tool/gin"
	"github.com/muxi-Infra/muxi-micro/tool/micro"
	"github.com/muxi-Infra/muxi-micro/tool/tests"
	"github.com/spf13/cobra"
)

func main() {
	var muxiCmd = &cobra.Command{
		Use:   "muxi",
		Short: "muxi-micro 总命令",
	}

	muxiCmd.AddCommand(curd.InitCurdCobra())
	muxiCmd.AddCommand(gin.InitGinCobra())
	muxiCmd.AddCommand(tests.InitTestCobra())
	muxiCmd.AddCommand(micro.InitMicroCobra())

	_ = muxiCmd.Execute()
}
