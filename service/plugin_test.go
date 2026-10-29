// plugin_test.go 覆盖服务注册校验与 go-plugin 进程内握手。
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/slaincsuite/modality-gopn-sdk-go/core"
)

type emptyServer interface {
	Ping()
}

type emptyServerImpl struct{}

func (emptyServerImpl) Ping() {}

var emptyServiceDesc = grpc.ServiceDesc{
	ServiceName: "test.Empty",
	HandlerType: (*emptyServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "test",
}

func TestPluginRegister(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr error
	}{
		{name: "empty name", service: NewService("", func(grpc.ServiceRegistrar) {}), wantErr: ErrInvalidService},
		{name: "nil register", service: NewService("test.Empty", nil), wantErr: ErrInvalidService},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPlugin().Register(tt.service)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Register() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginRegisterDuplicate(t *testing.T) {
	plugin := NewPlugin()
	service := NewService("test.Empty", func(registrar grpc.ServiceRegistrar) {
		registrar.RegisterService(&emptyServiceDesc, emptyServerImpl{})
	})
	if err := plugin.Register(service); err != nil {
		t.Fatalf("Register() = %v", err)
	}
	if err := plugin.Register(service); !errors.Is(err, ErrDuplicateService) {
		t.Fatalf("second Register() = %v, want ErrDuplicateService", err)
	}
}

// TestPluginServeWithoutRPC 零服务插件(仅应用/UI)仍须完成 go-plugin 握手并常驻。
func TestPluginServeWithoutRPC(t *testing.T) {
	serveCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reattachCh := make(chan *plugin.ReattachConfig, 1)
	serveErr := make(chan error, 1)
	pluginInstance := NewPlugin(
		WithLogger(core.NopLogger()),
		WithTestConfig(&plugin.ServeTestConfig{
			Context:          serveCtx,
			ReattachConfigCh: reattachCh,
			CloseCh:          make(chan struct{}),
		}),
	)
	go func() { serveErr <- pluginInstance.Serve() }()
	select {
	case <-reattachCh:
	case err := <-serveErr:
		t.Fatalf("Serve() 提前返回: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("等待握手配置超时")
	}
}

func TestPluginServeHandshake(t *testing.T) {
	serveCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reattachCh := make(chan *plugin.ReattachConfig, 1)
	closeCh := make(chan struct{})
	serveErr := make(chan error, 1)
	registered := make(chan struct{}, 1)
	pluginInstance := NewPlugin(
		WithLogger(core.NopLogger()),
		WithTestConfig(&plugin.ServeTestConfig{
			Context:          serveCtx,
			ReattachConfigCh: reattachCh,
			CloseCh:          closeCh,
		}),
	)
	err := pluginInstance.Register(NewService("test.Empty", func(registrar grpc.ServiceRegistrar) {
		registrar.RegisterService(&emptyServiceDesc, emptyServerImpl{})
		registered <- struct{}{}
	}))
	if err != nil {
		t.Fatalf("Register() = %v", err)
	}
	go func() {
		serveErr <- pluginInstance.Serve()
	}()
	var config *plugin.ReattachConfig
	select {
	case config = <-reattachCh:
	case err := <-serveErr:
		t.Fatalf("Serve() 提前返回: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("等待握手配置超时")
	}
	select {
	case <-registered:
	case <-time.After(3 * time.Second):
		t.Fatal("服务注册回调未执行")
	}
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  uint(core.ProtocolVersion),
			MagicCookieKey:   core.MagicCookieKey,
			MagicCookieValue: core.MagicCookieValue,
		},
		Reattach:         config,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclog.NewNullLogger(),
	})
	raw, err := client.Client()
	if err != nil {
		t.Fatalf("Client() = %v", err)
	}
	if err := raw.Ping(); err != nil {
		t.Fatalf("Ping() = %v", err)
	}
	cancel()
	select {
	case <-closeCh:
	case <-time.After(3 * time.Second):
		t.Fatal("等待插件服务退出超时")
	}
	if err := <-serveErr; err != nil {
		t.Fatalf("Serve() = %v", err)
	}
}
