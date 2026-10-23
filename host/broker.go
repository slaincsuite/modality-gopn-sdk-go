// broker.go 将 go-plugin broker 适配为宿主能力调用通道。
package host

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// BrokerServiceName 是宿主能力服务的固定 broker 服务标识。
const BrokerServiceName = "modality.v1.Host"

// BrokerCaller 依赖宿主在相同 broker ID 上提供 HostCapability 服务。
// 该实现只建立通道;最终是否放行仍由宿主逐次鉴权决定。
type BrokerCaller struct {
	broker *plugin.GRPCBroker
	conn   *grpc.ClientConn
}

// NewBrokerCaller 建立到宿主能力服务的 broker 连接。
func NewBrokerCaller(broker *plugin.GRPCBroker) (*BrokerCaller, error) {
	if broker == nil {
		return nil, ErrHostUnavailable
	}
	conn, err := broker.Dial(1)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHostUnavailable, err)
	}
	return &BrokerCaller{broker: broker, conn: conn}, nil
}

// Call 发起已授权的宿主访问;transport 为宿主生成的 gRPC 客户端。
func (c *BrokerCaller) Call(ctx context.Context, _ Authorization, _ Request) (Response, error) {
	if c == nil || c.conn == nil {
		return Response{}, ErrHostUnavailable
	}
	return Response{}, fmt.Errorf("host: 未提供 %s 客户端", BrokerServiceName)
}

// Close 关闭 broker 连接。
func (c *BrokerCaller) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
