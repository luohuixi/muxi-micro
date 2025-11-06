package create

import (
	"os"
	"path"

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

func CreateAllService(addr, dir string, apis []*parse.Api) error {
	for _, api := range apis {
		if err := CreateService(path.Join(addr, dir, "handler"), api); err != nil {
			return err
		}
	}
	if err := CreateRouter(path.Join(addr, dir, "router"), dir, apis); err != nil {
		return err
	}
	if err := CreateMain(path.Join(addr, dir), dir); err != nil {
		return err
	}
	if err := CreateWire(path.Join(addr, dir), dir, apis); err != nil {
		return err
	}
	return nil
}
