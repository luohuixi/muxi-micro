package create

import (
	"os"
	"path"
	"path/filepath"

	"github.com/muxi-Infra/muxi-micro/tool/gin/parse"
)

func CreateService(addr string, api *parse.Api) error {
	dir := path.Join(addr, api.ServiceName)
	// 存在不覆盖
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	if err := CreateType(dir, api.T, api.ServiceName); err != nil {
		return err
	}
	if err := Create2Service(dir, api.ServiceName, api.Service); err != nil {
		return err
	}
	if err := CreateHandler(dir, api.ServiceName, api.Service, api.Server); err != nil {
		return err
	}
	if err := CreateLogic(dir, api.ServiceName, api.Service); err != nil {
		return err
	}
	return nil
}

func CreateAllService(addr string, apis []*parse.Api) error {
	for _, api := range apis {
		if err := CreateService(path.Join(addr, "handler"), api); err != nil {
			return err
		}
	}
	// 获取根目录名
	dir := GetDirName(addr)
	if err := CreateRouter(path.Join(addr, "router"), dir, apis); err != nil {
		return err
	}
	if err := CreateMain(addr, dir); err != nil {
		return err
	}
	return nil
}

func GetDirName(addr string) string {
	currentDir, _ := os.Getwd()
	return filepath.Base(filepath.Join(currentDir, addr))
}
