// errors.go 定义 SDK 哨兵错误。
package service

import "errors"

var (
	// ErrDuplicateService 表示服务名重复。
	ErrDuplicateService = errors.New("modality: 服务名重复")
	// ErrInvalidService 表示服务名为空或注册函数为空。
	ErrInvalidService = errors.New("modality: 服务配置无效")
	// ErrHostClientUnsupported 表示插件进程不支持宿主侧 Dispense。
	ErrHostClientUnsupported = errors.New("modality: 插件进程不提供宿主客户端")
)
