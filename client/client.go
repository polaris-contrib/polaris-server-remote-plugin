/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package client

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/hashicorp/go-plugin"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/plugins"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/watcher"
)

// Client is a rich plugin client
type Client struct {
	sync.Mutex          // Mutex for manage client
	enable       bool   // represents the plugin is enabled or not
	pluginName   string // the name of the plugin, used for manage.
	pluginPath   string // the full path of the plugin, go-plugin start plugin according plugin path.
	address      string // address uses to interact with remote plugin server.
	watch        *watcher.Watcher
	config       *Config               // the setup config of the plugin
	pluginClient *plugin.Client        // the go-plugin client, polaris-serverImp run in grpc-client side.
	service      *plugins.PluginClient // service is the plugin-grpc-service client.
}

// Call invokes the function synchronously.
func (c *Client) Call(ctx context.Context, request *api.Request) (*api.Response, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	return c.service.Call(ctx, request)
}

// Close cleans up grpc connection.
func (c *Client) Close() error {
	if c.pluginClient == nil {
		return nil
	}
	if c.pluginClient.Exited() {
		return nil
	}
	if c.config.Mode == RumModelRemote {
		// remote plugin not support to shut down.
		return nil
	}
	c.pluginClient.Kill()
	return nil
}

// Disable set plugin is disabled.
func (c *Client) Disable() error {
	c.Lock()
	defer c.Unlock()

	if c.enable == false {
		return fmt.Errorf("plugin %s alread disbled", c.pluginName)
	}

	c.enable = false

	if c.pluginClient != nil {
		c.pluginClient.Kill()
	}
	return nil
}

// Open set plugin is enabled.
func (c *Client) Open() error {
	c.Lock()
	defer c.Unlock()

	c.enable = true
	return c.Check()
}

// Check checks client still alive, create if not alive
func (c *Client) Check() error {
	c.Lock()
	defer c.Unlock()

	if !c.enable {
		return errors.New("plugin " + c.pluginName + " disable")
	}

	// plugin still alive, return as early as possible.
	if c.pluginClient != nil && !c.pluginClient.Exited() {
		return nil
	}

	return c.recreate()
}

// newClient returns a new client
func newClient(config *Config) (*Client, error) {
	c := new(Client)
	c.enable = true
	c.pluginName = config.Name
	config.repairConfig()
	if config.Mode == RumModelRemote {
		c.address = config.Remote.Address
	} else {
		fullPath, err := config.pluginLoadPath()
		if err != nil {
			return nil, err
		}
		c.pluginPath = fullPath
		c.watch = watcher.New(fullPath, c.reload)
	}

	c.config = config
	if c.config.Logger == nil {
		c.config.Logger = log.DefaultLogger
	}
	if err := c.Check(); err != nil {
		return nil, fmt.Errorf("fail to check client: %w", err)
	}
	return c, nil
}

// recreate
func (c *Client) recreate() error {
	var (
		err          error
		pluginClient *plugin.Client
	)
	switch c.config.Mode {
	case RumModelLocal:
		pluginClient = c.recreateLocal()
	case RumModelRemote:
		pluginClient, err = c.recreateRemote()
	default:
		return fmt.Errorf("unkown plugin run mode: %s", c.config.Mode)
	}
	if err != nil {
		return err
	}

	if err = c.dispense(pluginClient); err != nil {
		return err
	}

	return nil
}

func (c *Client) dispense(pluginClient *plugin.Client) error {
	rpcClient, err := pluginClient.Client()
	if err != nil {
		return err
	}

	raw, err := rpcClient.Dispense(c.pluginName)
	if err != nil {
		return err
	}

	c.pluginClient = pluginClient
	c.service = raw.(*plugins.PluginClient)
	return nil
}

func (c *Client) recreateLocal() *plugin.Client {
	cmd := exec.Command(c.pluginPath, c.config.Local.Args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PLUGIN_PROCS=%d", c.config.Local.MaxProcs))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PLUGIN_NAME=%s", c.config.Name))
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: pluginsdk.Handshake,
		Plugins: map[string]plugin.Plugin{
			c.pluginName: &plugins.Plugin{},
		},
		Cmd:              cmd,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           c.config.Logger,
	})

	return pluginClient
}

func (c *Client) recreateRemote() (*plugin.Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", c.config.Remote.Address)
	if err != nil {
		return nil, err
	}
	pluginClient := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: pluginsdk.Handshake,
		Plugins: map[string]plugin.Plugin{
			c.pluginName: &plugins.Plugin{},
		},
		Reattach: &plugin.ReattachConfig{
			Protocol:        plugin.ProtocolGRPC,
			ProtocolVersion: int(pluginsdk.Handshake.ProtocolVersion),
			Addr:            addr,
			// Mock plugin process can be found, but we must note that
			// this method may be unstable.
			Pid: os.Getpid(),
		},
		Logger:           c.config.Logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	return pluginClient, nil
}

func (c *Client) reload(_ string) {
	if c == nil || c.pluginClient == nil {
		return
	}
	if err := c.recreate(); err != nil {
		log.DefaultLogger.Fatal("recreate error", "error", err)
	}
	log.DefaultLogger.Info("plugin recreate finished")
}
