package service

import (
	"git.ezbuy.me/ezbuy/evarmor/common/metcd"
)

// Beacon 烽火台
func Beacon(addr []string) {
	cli, _ := metcd.NewClientDis([]string{"127.0.0.1:2379"})

	go func() {
		if results, err := cli.GetService("/handler"); err == nil {
			beaconServiceHandler(results)
		}
	}()
}

// 监听注册的接口服务
func beaconServiceHandler(handles []string) error {
	// 写入代理服务
	return nil
}
