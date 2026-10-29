// manifest/errors.go 定义插件清单校验的哨兵错误。
package manifest

import "errors"

var (
	// ErrInvalidManifest 表示清单文件不存在、无法解析或包含非法字段。
	ErrInvalidManifest = errors.New("manifest: 清单无效")
	// ErrInvalidActions 表示操作目录(action.json)不存在、无法解析或包含非法字段。
	ErrInvalidActions = errors.New("manifest: 操作目录无效")
	// ErrHostVersionIncompatible 表示插件声明的最低宿主版本与目标宿主不兼容。
	ErrHostVersionIncompatible = errors.New("manifest: 宿主版本不兼容")
	// ErrDuplicateName 表示扫描到的插件目录存在重复名称。
	ErrDuplicateName = errors.New("manifest: 插件名重复")
)
