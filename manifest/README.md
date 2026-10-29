# manifest

插件包清单工具:读取并校验 `plugin.json`、插件二进制、UI 入口、
`routes.json` 与操作目录 `action.json`。校验规则与 modality 宿主
`library/gopn/manifest` 对齐,便于插件作者在打包前发现布局问题。

```go
manifest, err := manifest.Load("./build/plugins/my-plugin")
if err != nil {
    log.Fatal(err)
}
fmt.Println(manifest.Name, manifest.Version)
```

`ValidationError` 会汇总全部字段错误,可经 `errors.As` 获取字段路径。

## 操作目录(action.json)

`action.json` 与 `plugin.json` 同级存放,声明插件对宿主暴露的应用
(策略的「服务」)、操作(策略的「动作」)与数据源;同一份声明也是插件前端
权限的单一来源,插件作者不得在插件前端另抄一份操作清单。
`plugin.json` 的 `actions` 字段可指向其他路径,缺省自动探测根目录 `action.json`;
两者都不存在时为纯服务插件,`Manifest.Apps` 保持为空。

- 应用编码与操作名 kebab-case 且唯一,合并为 sam 功能名 `应用编码.操作名`;
- `kind` 取 `read` 读 / `write` 写 / `list` 列表 / `other` 其他四类;
- 操作的 `resource` 必须引用应用已声明的 `resources[].code`。

```go
apps, err := manifest.LoadActions("./action.json")
for _, app := range apps {
    for _, action := range app.Actions {
        fmt.Println(app.Code+"."+action.Name, action.Kind)
    }
}
```

`Load` 会自动加载 `action.json` 并把结果填充到 `Manifest.Apps`;
解析或校验失败时错误同时满足 `errors.Is(err, ErrInvalidManifest)` 与
`errors.Is(err, ErrInvalidActions)`。

## API

| 函数 / 方法 | 用途 |
| --- | --- |
| `Load(dir)` | 读取并校验单个插件目录(含 `action.json`) |
| `Discover(dir)` | 扫描直接子目录,按名称排序;重复名报错 |
| `CheckHostVersion(hostVersion)` | 校验 `min_host_version` 范围 |
| `LoadRoutes(path)` | 解析 `routes.json` |
| `LoadActions(path)` | 解析并校验 `action.json`,返回应用声明 |
| `ValidOperationKind(kind)` | 判定操作类别是否为四类之一 |

| 类型 / 常量 | 用途 |
| --- | --- |
| `App` / `Action` / `Resource` | 应用、操作与数据源的声明模型 |
| `OpKindRead` / `OpKindWrite` / `OpKindList` / `OpKindOther` | 四类操作枚举 |
| `ActionFileName` | 操作目录缺省文件名 `action.json` |
