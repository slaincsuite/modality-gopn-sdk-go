# modality Plugin SDK

modality Go 插件 SDK。插件是独立二进制进程,由 modality 宿主通过
gopn 启动;可暴露 gRPC 服务,也可只提供**应用与操作目录(action.json)/UI**
(零服务插件合法,仍以 go-plugin 握手与内建健康服务常驻)。SDK 提供进程握手、
服务注册、日志桥接、清单校验、扩展点声明与默认拒绝的宿主能力调用。

> 宿主侧能力通道(broker)尚未实现:`grants` 声明暂不会被授予或调用。

## 环境要求

- Go 1.27.3+
- `protoc`、`protoc-gen-go`、`protoc-gen-go-grpc`(仅生成业务 proto 时需要)

## 安装

```sh
go get github.com/slaincsuite/modality-gopn-sdk-go@v0.1.3
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

插件包以 `.mgpn` 压缩包交付,内含以下目录结构:

```text
demo-plugin/
├── plugin.json          # 清单(项目根目录)
├── action.json          # 操作目录:应用、操作与数据源声明(可选)
├── src/                 # Go 源码
│   ├── main.go          # 进程入口(src/apps/main.go)
│   └── starter.go       # 业务实现
└── ui/                  # 前端工程(可选)
    ├── index.html       # SPA 入口
    ├── src/
    │   └── routes.json  # 前端路由声明
    └── assets/          # 构建产物
```

`plugin.json` 位于项目根目录;
`ui/routes_file` 引用 `ui/src/routes.json`。
`plugin.json` 的 `services` 必须与 `service.Service.Name` 完全一致;
`grants` 只是能力申请,不代表宿主授权结果。发布前使用 `manifest.Load`
校验目录。

## 应用与操作声明

`action.json` 声明插件对宿主暴露的应用、操作与数据源,是宿主策略
(「服务」= 应用编码,「操作」= 操作名)与插件前端权限的单一来源:

```json
{
  "apps": [
    {
      "code": "demo-plugin",
      "name": "示例应用",
      "actions": [
        { "name": "record-list", "title": "查询记录列表", "kind": "list", "resource": "record" }
      ],
      "resources": [
        { "code": "record", "name": "记录" }
      ]
    }
  ]
}
```

- 应用编码与操作名 kebab-case,合并为宿主功能名 `应用编码.操作名`;
- `kind` 取 `read` / `write` / `list` / `other` 四类
  (插件前端按类约定组织界面:读=打开/查看,写=编辑,列表=看列表/导出,其他=其他入口);
- 插件前端只从宿主 `GET /open-api/v1/upms/policies/my-actions?app=<应用编码>`
  取权限,不得内置权限判断或另抄操作清单;
- `plugin.json` 的 `actions` 字段可指向其他路径,缺省自动探测根目录
  `action.json`;打包前用 `manifest.LoadActions` 在 Go 侧校验。

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
