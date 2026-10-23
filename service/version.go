// version.go 暴露 SDK 版本。
package service

import "github.com/slaincsuite/modality-gopn-sdk-go/core"

// Version 返回当前 SDK 版本。
func Version() string {
	return core.SDKVersion
}
