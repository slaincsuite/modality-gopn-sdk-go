# modality Plugin SDK

modality Go 插件 SDK。插件是独立二进制进程,由 modality 宿主通过
gopn 启动,并以 gRPC 暴露 proto 生成的服务。SDK 提供进程握手、服务注册、
日志桥接、清单校验、扩展点声明与默认拒绝的宿主能力调用。

## 环境要求

- Go 1.27.3+
- `protoc`、`protoc-gen-go`、`protoc-gen-go-grpc`(仅生成业务 proto 时需要)

## 安装

```sh
go get github.com/slaincsuite/modality-gopn-sdk-go
```

## 最小插件

```go
package main

import (
    "context"

    "google.golang.org/grpc"
    healthpb "google.golang.org/grpc/health/grpc_health_v1"

    "github.com/slaincsuite/modality-gopn-sdk-go/service"
)

type healthServer struct {
    healthpb.UnimplementedHealthServer
}

func (healthServer) Check(context.Context, *healthpb.HealthCheckRequest) (
    *healthpb.HealthCheckResponse, error,
) {
    return &healthpb.HealthCheckResponse{
        Status: healthpb.HealthCheckResponse_SERVING,
    }, nil
}

func main() {
    service.MustServe([]service.Service{
        service.NewService("starter.v1.Health", func(r grpc.ServiceRegistrar) {
            healthpb.RegisterHealthServer(r, healthServer{})
        }),
    })
}
```

## 宿主访问

宿主访问必须声明能力并默认拒绝:

```go
caller := // 由宿主反向通道注入
client := service.NewHostClient(service.HostConfig{
    PluginName: "demo-plugin",
    Capabilities: []service.Capability{{
        Name: "module.registry.read", Module: "module", Mode: service.AccessRead,
    }},
}, caller)

response, err := client.Call(ctx, "module.registry.read", nil)
```

`Authorization` 只是调用元数据;宿主仍会在每次调用时重新鉴权并审计。
受控写入必须由宿主显式开放,并逐次校验。

## 插件包

```text
demo-plugin/
├── plugin.json
├── demo-plugin
└── ui/
```

`plugin.json` 的 `services` 必须与 `service.Service.Name` 完全一致;
`grants` 只是能力申请,不代表宿主授权结果。发布前使用 `manifest.Load`
校验目录。

## 包索引

| 包 | 职责 |
| --- | --- |
| [service](service/README.md) | 插件组装、服务注册、Serve |
| [core](core/README.md) | 握手契约、配置、日志、版本 |
| [manifest](manifest/README.md) | 插件清单与路由校验 |
| [extensionpoint](extensionpoint/README.md) | 扩展点声明模型 |
| [host](host/README.md) | 宿主能力授权与调用 |

## 验证

```sh
make verify
```
