package create

import (
	"github.com/muxi-Infra/muxi-micro/tool/micro/parse"
)

// CreateServerScaffold 创建 handler 和 main 部分代码
func CreateServerScaffold(cover bool, output string, protoFile *parse.ProtoFile) error {
	if err := CreateHandler(cover, output, protoFile); err != nil {
		return err
	}

	if err := CreateRpc(cover, output, protoFile); err != nil {
		return err
	}

	if err := CreateMain(cover, output, protoFile); err != nil {
		return err
	}
	return nil
}
