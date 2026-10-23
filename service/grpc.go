// grpc.go 将插件服务桥接为 go-plugin 的 GRPCPlugin。
package service

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const pluginName = "modality"

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	services []Service
	broker   func(*plugin.GRPCBroker) error
}

// GRPCServer 注册全部服务;单个服务 panic 转为错误返回给 go-plugin。
func (p *grpcPlugin) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) (err error) {
	for _, service := range p.services {
		if registerErr := registerService(server, service); registerErr != nil {
			return registerErr
		}
	}
	if p.broker == nil {
		return nil
	}
	if err = p.broker(broker); err != nil {
		return fmt.Errorf("modality: broker: %w", err)
	}
	return nil
}

// GRPCClient 明确不支持宿主侧 Dispense。
func (p *grpcPlugin) GRPCClient(
	context.Context,
	*plugin.GRPCBroker,
	*grpc.ClientConn,
) (any, error) {
	return nil, ErrHostClientUnsupported
}

func registerService(server *grpc.Server, service Service) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("modality: %s 注册失败: %v", service.Name, recovered)
		}
	}()
	service.Register(server)
	return nil
}
