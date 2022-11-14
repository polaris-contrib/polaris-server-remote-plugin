package main

import (
	"context"
	"time"

	pluginsdk "github.com/polaris-contrib/polaris-server-remote-plugin-common"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/client"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
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
			Mode: client.RumModelCompanion,
			Companion: client.CompanionConfig{
				MaxProcs: 1,
			},
		},
	); err != nil {
		log.Fatal("server-v1 register failed", "error", err)
		return
	}

	if client2, err = client.Register(
		&client.Config{
			Name: "remote-rate-limit-server-v2",
			Mode: client.RumModelCompanion,
		},
	); err != nil {
		log.Fatal("server-v2 register failed", "error", err)
	}

	if client3, err = client.Register(
		&client.Config{
			Name: "remote-rate-limit-server-v3",
			Mode: client.RumModelRemote,
			Remote: client.RemoteConfig{
				Address: "0.0.0.0:8972",
			},
		}); err != nil {
		log.Fatal("server-v3 register failed", "error", err)
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
		log.Fatal("fail to marshal request", "client_name", name, "error", err)
	}

	response, err := pc.Call(context.Background(), req)

	if err != nil {
		log.Fatal("fail to invoke", "client_name", name, "error", err)
	}

	log.Info("response body from client", "client_name", name, "response", response.String())
}
