package registry

import (
	"context"
	"fmt"
	"net"
)

type RegistrationCenter interface {
	Register(ctx context.Context, serviceName, host, port string) error
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// 过滤掉 loopback 地址
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ip4 := ipNet.IP.To4(); ip4 != nil {
				return ip4.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no valid IP address found")
}
