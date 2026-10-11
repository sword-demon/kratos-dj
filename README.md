# kratos 代驾微服务项目


## 顾客获取验证码

```bash
kratos proto add api/customer/customer.proto
```

```bash
# 生成客户端代码
# 生成的代码位于 api/customer/
kratos proto client api/customer/customer.proto
# 生成服务端代码
kratos proto server api/customer/customer.proto
# 代码位于 internal/service/customer 目录下
```

### 在 http 服务中加入 customer 服务

编辑 `/internal/server/http.go`

```go
package server

import (
	"customer/api/customer"
	v1 "customer/api/helloworld/v1"
	"customer/internal/conf"
	"customer/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	customerService *service.CustomerService,
	greeter *service.GreeterService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	// 注册 customer 的 http 服务
	customer.RegisterCustomerHTTPServer(srv, customerService)
	// 注册 greeter 的 http 服务
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
```

### 处理 go 的相关依赖

更新 wire 的 Provider 定义

```go
var ProviderSet = wire.NewSet(NewCustomerService, NewGreeterService)
```

基于 wire 生成代码

```bash
go generate ./...
wire: customer/cmd/customer: wrote /Volumes/MOVESPEED/ai-coding/kratos-dj/backend/customer/cmd/customer/wire_gen.go
```

### 运行测试
```bash
kratos run
```

get 请求调用

```http request
http://localhost:8100/customer/get-verify-code
```

```json
{
  "Code": 0,
  "Message": "",
  "Data": ""
}
```