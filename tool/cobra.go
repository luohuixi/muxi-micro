package main

import (
	"github.com/muxi-Infra/muxi-micro/tool/curd"
	"github.com/spf13/cobra"
)

func main() {
	var muxiCmd = &cobra.Command{
		Use:   "muxi",
		Short: "muxi-micro 总命令",
	}

	curdCmd := curd.InitCurdCobra()
	muxiCmd.AddCommand(curdCmd)

	_ = muxiCmd.Execute()
}
