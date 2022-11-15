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

package plugins

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
)

// RPCClient wraps a remote plugin api client.
type RPCClient struct {
	api.PluginClient
}

// Service must implement plugin.PluginServer
var _ api.PluginServer = (Service)(nil)

// Service implements remote proto defined Plugin Server
//
//go:generate mockgen -source=plugin.go -destination=plugin_mock.go -package=plugins Service
type Service interface {
	// Call implement plugin.Server.Call
	Call(ctx context.Context, request *api.Request) (*api.Response, error)
}

// serverImp must implements Service.
var _ Service = (*serverImp)(nil)

// serverImp plugin serverImp.
type serverImp struct {
	Backend Service
}

// Call calls plugin backend serverImp.
func (s *serverImp) Call(ctx context.Context, req *api.Request) (*api.Response, error) {
	return s.Backend.Call(ctx, req)
}

// Plugin must implement plugin.GRPCPlugin
var _ plugin.GRPCPlugin = (*Plugin)(nil)

// Plugin is the implementation of plugin.GRPCPlugin, so we can serve/consume this.
type Plugin struct {
	plugin.Plugin
	Backend Service
}

// NewPlugin returns a new Plugin.
func NewPlugin(backend Service) *Plugin {
	return &Plugin{Backend: backend}
}

// GRPCServer implements plugin.Plugin GRPCServer method.
func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	api.RegisterPluginServer(s, &serverImp{Backend: p.Backend})
	return nil
}

// GRPCClient implements plugin.Plugin GRPCClient method.
func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &RPCClient{PluginClient: api.NewPluginClient(c)}, nil
}
