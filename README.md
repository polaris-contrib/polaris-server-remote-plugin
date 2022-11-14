# polaris-server-remote-plugin-common

## 实现方案

Polaris-Server 支持三种类型的插件机制：本地（local）、伴生（companion）、远端（remote）。基于 [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) 实现了伴生和远端插件的能力。本地插件的实现方式是基于插件和北极星 Server 运行在同一进程内，这要求使用Polaris 的用户需要基于源代码编译Polaris 服务，具体实现方式可以参考 [服务端插件开发](https://polarismesh.cn/docs/%E5%8F%82%E8%80%83%E6%96%87%E6%A1%A3/%E5%BC%80%E5%8F%91%E8%80%85%E6%96%87%E6%A1%A3/%E6%8F%92%E4%BB%B6%E5%BC%80%E5%8F%91/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8F%92%E4%BB%B6%E5%BC%80%E5%8F%91/) 文档。

 [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) 基于 gRPC 协议实现插件机制，主服务端通过 gRPC 的 Client 端访问插件服务的 Server 端，插件服务的 Server 端通过`exec.Cmd` 在主服务端进程中启动子进程拉起 gRPC 服务，主服务和插件服务通过本地网络栈实现网路通信。基于 [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) 的实现原理，拓展伴生、远端两种插件能力。

## companion 伴生插件

为了区别通过源代码编译实现的 [服务端插件开发](https://polarismesh.cn/docs/%E5%8F%82%E8%80%83%E6%96%87%E6%A1%A3/%E5%BC%80%E5%8F%91%E8%80%85%E6%96%87%E6%A1%A3/%E6%8F%92%E4%BB%B6%E5%BC%80%E5%8F%91/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E6%8F%92%E4%BB%B6%E5%BC%80%E5%8F%91/) 方案，将主服务通过二进制文件拉起 gRPC 插件服务的方式成为伴生模式，这种模式也就是 [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) 的原生方案。以拓展 [polaris server](https://github.com/polarismesh/polaris]) 的限流插件为例：

```go
const PluginName = "remote-rate-limit"
// 插件注册
func init() {
	plugin.RegisterPlugin(PluginName, &RateLimiter{})
}

// RateLimiter 远程限流插件
type RateLimiter struct {
	cfg    *client.Config
	client *client.Client
}

```

### 配置

注册名为 `remote-rate-limit` 的远程插件，通过 mode 指定为插件模式为`companion`，即伴生模式。由于伴生模式需要加载二进制文件，因此需要指定文件的路径：

- 可以从通过从主进程工作目录查找名为 `name` 字段指定的二进制文件；
- 可以通过 `companion.path` 字段指定二进制文件的相对/绝对路径；

```yaml
plugin:
  ratelimit:
    name: remote-rate-limit
    option:
      name: rate-limit-server-v2
      mode: companion
      companion:
        path: ../polaris-server-remote-plugin-common/examples/remote-rate-limit-server-v2
        max-procs: 4
```

同时，也可以通过 `companion.max-procs` 指定插件服务进程可以使用的最大 M 数量。



### 注册

```go
// Initialize 初始化函数
func (r *RateLimiter) Initialize(c *plugin.ConfigEntry) error {
	... // 配置读取 
	r.client, err = client.Register(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup rate-limit plugin: %w", err)
	}
}
```



### 销毁接口

```go
// Destroy 销毁函数
func (r *RateLimiter) Destroy() error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}
```



### 插件服务开发

```go
package main

import (
	"context"

	pluginsdk "github.com/polaris-contrib/polaris-server-remote-plugin-common"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/server"
)

type rateLimiter struct{}

func (s *rateLimiter) Call(_ context.Context, request *api.Request) (*api.Response, error) {
	var rateLimitRequest api.RateLimitPluginRequest
	if err := pluginsdk.UnmarshalRequest(request, &rateLimitRequest); err != nil {
		return nil, err
	}
	response, err := pluginsdk.MarshalResponse(&api.RateLimitPluginResponse{Allow: true})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func main() {
	server.Serve(&rateLimiter{})
}
```







## Remote 远程插件



### 配置

```yaml
plugin:
  ratelimit:
    name: remote-rate-limit
    option:
      name: rate-limit-server-v3
      mode: remote
      remote:
        address: 0.0.0.0:8972
```



### 插件服务开发

```go
type rateLimiter struct {
}

func (r rateLimiter) Call(_ context.Context, request *api.Request) (*api.Response, error) {
	var rateLimitRequest api.RateLimitPluginRequest
	if err := pluginsdk.UnmarshalRequest(request, &rateLimitRequest); err != nil {
		return nil, err
	}
	allow := false
	n, _ := rand.Int(rand.Reader, big.NewInt(3))
	if n.Int64() == 2 {
		allow = true
	}

	response, err := pluginsdk.MarshalResponse(&api.RateLimitPluginResponse{Allow: allow})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                       // 创建gRPC服务器
	api.RegisterPluginServer(s, &rateLimiter{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
```

