# Polaris Server Remote Plugin

This is a common library for the Polaris Server Remote Plugin. You can use this library to create your own remote
plugin. This package allows you to run you own grpc-based plugin service, you just need to implement the
api.PluginServer interface.

In order to ensure the stability and backward compatibility of the `plugin proto`, we use the Any type. Your service
only needs to implement the development of the plug -in to complete the plug -in of each interface of the ping and call
chain.

## How to write a plugin ?

You can refer to the example code, and you can use any programming language to write the plugin services.

```go
type rateLimiter struct{}

func (s *rateLimiter) Ping(_ context.Context, request *api.PingRequest) (*api.PongResponse, error) {
    log.Info("ping pong")
    return &api.PongResponse{}, nil
}

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
```

## How and where can we find your plugin ?

As shown in the above code, you have written the code for your service. We agreed to read all the Unix Domain Socket
files in the `${POLARIS_PLUGGABLE_SOCKETS_FOLDER}` directory, and try to establish a connection with your service
through the socket file. In order to
ensure that your service meets our plugin services, we need you to turn on the GRPC reflection service.Therefore, your
service needs to be exposed through the socket, and the reflection service needs to be turned on.

```go
func main() {
	sockFile := "/tmp/polaris-pluggable-sockets/rate-limit.sock"
	cleanup(sockFile) // clean up the old socket file

	lis, err := net.Listen("unix", sockFile)
	if err != nil {
		log.Fatal("failed to listen: %+v", err)
	}

	s := grpc.NewServer()
	api.RegisterPluginServer(s, &rateLimiter{})

	reflection.Register(s) // enable reflection service

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
```
