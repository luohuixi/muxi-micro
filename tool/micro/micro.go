package micro

import (
	"github.com/muxi-Infra/muxi-micro/tool/micro/create"
	"github.com/muxi-Infra/muxi-micro/tool/micro/parse"
	"github.com/spf13/cobra"
)

func InitMicroCobra() *cobra.Command {
	var microCmd = &cobra.Command{
		Use:   "micro",
		Short: "micro 生成 grpc 微服务框架",
	}

	microCmd.RunE = makeRunE(func(cover bool, output string, file *parse.ProtoFile) error {
		if err := create.CreateServer(cover, output, file); err != nil {
			return err
		}
		if err := create.CreateClient(cover, output, file); err != nil {
			return err
		}
		if err := create.CreateServerScaffold(cover, output, file); err != nil {
			return err
		}
		return nil
	})

	microCmd.PersistentFlags().String("proto", "./hello.proto", "proto 文件的位置")
	microCmd.PersistentFlags().String("output", ".", "输出路径")
	microCmd.PersistentFlags().Bool("skip-protoc", false, "是否跳过执行 protoc 命令")
	microCmd.PersistentFlags().Bool("cover", false, "是否覆盖原文件")

	var ServerCmd = &cobra.Command{
		Use:   "server",
		Short: "生成封装好的 grpc server",
		RunE: makeRunE(func(cover bool, output string, file *parse.ProtoFile) error {
			return create.CreateServer(cover, output, file)
		}),
	}

	var ClientCmd = &cobra.Command{
		Use:   "client",
		Short: "生成封装好的 grpc client",
		RunE: makeRunE(func(cover bool, output string, file *parse.ProtoFile) error {
			return create.CreateClient(cover, output, file)
		}),
	}

	microCmd.AddCommand(ServerCmd)
	microCmd.AddCommand(ClientCmd)

	return microCmd
}

func makeRunE(handler func(cover bool, output string, file *parse.ProtoFile) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		proto, _ := cmd.Flags().GetString("proto")
		output, _ := cmd.Flags().GetString("output")
		skip, _ := cmd.Flags().GetBool("skip-protoc")
		cover, _ := cmd.Flags().GetBool("cover")

		file, err := parse.ParseProto(proto)
		if err != nil {
			return err
		}

		if err := create.GenerateProto(proto, output, skip); err != nil {
			return err
		}
		return handler(cover, output, file)
	}
}
