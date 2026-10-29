// serve.go 提供一条调用完成注册与运行的快捷入口。
package service

import (
	"fmt"
	"os"
)

// Serve 创建插件、注册全部服务并阻塞运行。
func Serve(services []Service, options ...Option) error {
	instance := NewPlugin(options...)
	for _, service := range services {
		if err := instance.Register(service); err != nil {
			return err
		}
	}
	return instance.Serve()
}

// MustServe 运行插件;配置或运行失败时输出到 stderr 并以非零状态退出。
func MustServe(services []Service, options ...Option) {
	if err := Serve(services, options...); err != nil {
		fmt.Fprintf(os.Stderr, "modality plugin: %v\n", err)
		os.Exit(1)
	}
}
