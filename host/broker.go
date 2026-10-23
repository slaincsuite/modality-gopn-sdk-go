// broker.go 将 go-plugin broker 适配为宿主能力调用通道。
package host

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// BrokerServiceName 是宿主能力服务的固定 broker 服务标识。
const BrokerServiceName = "modality.v1.Host"

// BrokerCaller 依赖宿主在相同 broker ID 上提供 modality.v1.Host 服务。
// 该实现只建立通道并转发调用;最终是否放行仍由宿主逐次鉴权决定。
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

// Call 发起已授权的宿主访问。载荷经 structpb 序列化,以 JSON 兼容格式
// 通过 gRPC 通用通道发送;宿主收到后重新鉴权并按能力契约返回数据。
func (c *BrokerCaller) Call(ctx context.Context, _ Authorization, req Request) (Response, error) {
	if c == nil || c.conn == nil {
		return Response{}, ErrHostUnavailable
	}
	payload, err := structpb.NewValue(req.Payload)
	if err != nil {
		return Response{}, fmt.Errorf("%w: 载荷序列化: %v", ErrInvalidCapability, err)
	}
	inVal := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"capability": structpb.NewStringValue(req.Capability),
			"payload":    payload,
		},
	}
	inBytes, err := proto.Marshal(inVal)
	if err != nil {
		return Response{}, fmt.Errorf("%w: proto 编码: %v", ErrInvalidCapability, err)
	}
	method := "/" + BrokerServiceName + "/Call"
	var outBytes []byte
	if err := c.conn.Invoke(ctx, method, inBytes, &outBytes); err != nil {
		return Response{}, fmt.Errorf("host: 能力调用 %s: %w", req.Capability, err)
	}
	var outVal structpb.Struct
	if err := proto.Unmarshal(outBytes, &outVal); err != nil {
		return Response{}, fmt.Errorf("host: 响应解码: %w", err)
	}
	dataField, ok := outVal.Fields["data"]
	if !ok {
		return Response{}, nil
	}
	return Response{Data: dataField.AsInterface()}, nil
}

// Close 关闭 broker 连接。
func (c *BrokerCaller) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
