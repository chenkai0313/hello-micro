# hello-micro
基于go-micro/v2 封装的微服务框架

### 使用
```
go get github.com/chenkai0313/hello-micro
```

#### 使用事例
go build main.go

./main


#### 说明
- 系统使用的go-micro/v2版本
- 默认封装了jaeger (做了改动，配合使用gin使用)，初次使用，可以注释掉
- 封装了基于zaplog 日志管理工具,本人是结合elk进行使用的，需要存储数据库的自行修改
- 封装了基于gorm/v2的
- 封装了自定义的validator
- 生产/测试环境 选择基于环境变量"env"
- 内涵盖了常用的的方法(unitls目录下)
- 封装了图片验证码
- 封装了rabbitMq(包含了延迟队列方法封装)
- 封装了redis
- 基于makefile 的protobuf
- 服务发现使用的是etcd
- 更多内容自行挖掘
- 全部个人封装,仅供学习参考
