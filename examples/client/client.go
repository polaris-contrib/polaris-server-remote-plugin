package main

import (
	"context"
	"log"
	"time"

	pluginsdk "github.com/polaris-contrib/polaris-server-remote-plugin-common"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/client"
)

func main() {
	var (
		err     error
		client1 *client.Client
		client2 *client.Client
		client3 *client.Client
	)
	if client1, err = client.Register(
		&client.Config{
			Name: "remote-rate-limit-server-v1",
			Mode: client.RumModelLocal,
			Local: client.LocalConfig{
				MaxProcs: 1,
			},
		},
	); err != nil {
		log.Fatalf("server-v1 register failed: %+v", err)
		return
	}

	if client2, err = client.Register(
		&client.Config{
			Name: "remote-rate-limit-server-v2",
			Mode: client.RumModelLocal,
		},
	); err != nil {
		log.Fatalf("server-v2 register failed: %+v", err)
	}

	if client3, err = client.Register(
		&client.Config{
			Name: "remote-rate-limit-server-v3",
			Mode: client.RumModelRemote,
			Remote: client.RemoteConfig{
				Address: "0.0.0.0:8972",
			},
		}); err != nil {
		log.Fatalf("server-v3 register failed: %+v", err)
	}

	clientInvoke(client1, "1")
	clientInvoke(client2, "2")

	for i := 0; i < 1000; i++ {
		clientInvoke(client3, "3")
		time.Sleep(time.Second)
	}
}

func clientInvoke(pc *client.Client, name string) {
	req, err := pluginsdk.MarshalRequest(&api.RateLimitPluginRequest{
		Type: api.RatelimitType_IPRatelimit, Key: "127.0.0.1",
	})
	if err != nil {
		log.Fatalf("client-%s fail to marshal request: %+v", name, err)
	}

	response, err := pc.Call(context.Background(), req)

	if err != nil {
		log.Fatalf("client-%s fail to invoke: %+v", name, err)
	}

	log.Printf("response body from client-%s: %s\n", name, response.String())
}
