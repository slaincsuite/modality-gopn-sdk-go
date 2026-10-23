# core

SDK 基础包:与 modality 宿主 `library/gopn` 对齐的握手常量、配置、日志接口
与默认实现。

## 关键契约

| 项 | 值 |
| --- | --- |
| 协议版本 | `1` |
| Magic Cookie Key | `GOPN_PLUGIN` |
| Magic Cookie Value | `gopn-plugin-handshake` |

前三个值必须与宿主保持一致。插件作者通常只需使用根包
`modality.WithLogger` / `modality.WithLogLevel`,无需直接构造 `core.Config`。
根包同时重新导出 `Logger`、`LogLevel` 与日志级别常量,一般无需直接导入
`core`。
