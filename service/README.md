# service

插件进程组装门面:提供 `Plugin`、`Service`、`Serve`、日志选项与公共类型
重导出。插件作者在这里注册 proto 生成的 gRPC 服务并阻塞等待宿主调用。
