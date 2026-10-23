// core/handshake.go 定义与 modality 宿主 gopn 对齐的 go-plugin 握手契约。
package core

import (
	"errors"
	"fmt"
)

const (
	// ProtocolVersion 当前插件握手协议版本。
	ProtocolVersion = 1
	// MagicCookieKey 与 modality 宿主 library/gopn 保持一致的握手键。
	MagicCookieKey = "GOPN_PLUGIN"
	// MagicCookieValue 与 modality 宿主 library/gopn 保持一致的握手值。
	MagicCookieValue = "gopn-plugin-handshake"
)

// ErrInvalidHandshake 表示握手配置不完整或协议版本非法。
var ErrInvalidHandshake = errors.New("core: 握手配置无效")

// Handshake 描述宿主与插件共同使用的 go-plugin 握手参数。
type Handshake struct {
	ProtocolVersion  int
	MagicCookieKey   string
	MagicCookieValue string
}

// DefaultHandshake 返回当前 modality 宿主的握手参数。
func DefaultHandshake() Handshake {
	return Handshake{
		ProtocolVersion:  ProtocolVersion,
		MagicCookieKey:   MagicCookieKey,
		MagicCookieValue: MagicCookieValue,
	}
}

// Validate 校验握手参数;ProtocolVersion 必须为正数。
func (h Handshake) Validate() error {
	if h.ProtocolVersion <= 0 {
		return fmt.Errorf("%w: ProtocolVersion 必须大于 0", ErrInvalidHandshake)
	}
	if h.MagicCookieKey == "" {
		return fmt.Errorf("%w: MagicCookieKey 必填", ErrInvalidHandshake)
	}
	if h.MagicCookieValue == "" {
		return fmt.Errorf("%w: MagicCookieValue 必填", ErrInvalidHandshake)
	}
	return nil
}
