// errors.go 定义宿主访问与授权哨兵错误。
package host

import "errors"

var (
	// ErrForbidden 表示授权不覆盖目标模块、操作或条件。
	ErrForbidden = errors.New("host: 插件无权访问")
	// ErrInvalidCapability 表示能力声明不完整或模式非法。
	ErrInvalidCapability = errors.New("host: 能力声明无效")
	// ErrHostUnavailable 表示宿主通道不可用。
	ErrHostUnavailable = errors.New("host: 宿主不可用")
)
