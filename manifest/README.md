# manifest

插件包清单工具:读取并校验 `plugin.json`、插件二进制、UI 入口与
`routes.json`。校验规则与 modality 宿主 `library/gopn/manifest` 对齐,
便于插件作者在打包前发现布局问题。

```go
manifest, err := manifest.Load("./build/plugins/my-plugin")
if err != nil {
    log.Fatal(err)
}
fmt.Println(manifest.Name, manifest.Version)
```

`ValidationError` 会汇总全部字段错误,可经 `errors.As` 获取字段路径。

## API

| 函数 / 方法 | 用途 |
| --- | --- |
| `Load(dir)` | 读取并校验单个插件目录 |
| `Discover(dir)` | 扫描直接子目录,按名称排序;重复名报错 |
| `CheckHostVersion(hostVersion)` | 校验 `min_host_version` 范围 |
| `LoadRoutes(path)` | 解析 `routes.json` |
