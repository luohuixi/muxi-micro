package main

import (
	"github.com/muxi-Infra/muxi-micro/tool/curd"
	"github.com/muxi-Infra/muxi-micro/tool/new"
	"github.com/spf13/cobra"
)

func main() {
	var muxiCmd = &cobra.Command{
		Use:   "muxi",
		Short: "muxi-micro 总命令",
	}

	muxiCmd.AddCommand(curd.InitCurdCobra())
	muxiCmd.AddCommand(new.InitNewCobra())

	_ = muxiCmd.Execute()
}
