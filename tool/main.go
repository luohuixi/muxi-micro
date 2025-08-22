package main

import (
	curd "github.com/muxi-Infra/muxi-micro/tool/curd/cobra"
	"github.com/spf13/cobra"
)

func main() {
	var Mu_xiCmd = &cobra.Command{
		Use:   "Mu_xi",
		Short: "Mu_xi-micro 总命令",
	}

	curd.InitCurdCobra(Mu_xiCmd)

	_ = Mu_xiCmd.Execute()
}
