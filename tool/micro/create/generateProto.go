package create

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func GenerateProto(proto, output string, skip bool) error {
	if skip {
		return nil
	}

	cmd := exec.Command("protoc",
		fmt.Sprintf("--go_out=%s", output),
		fmt.Sprintf("--go-grpc_out=%s", output),
		proto)

	// 捕获输出
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}
