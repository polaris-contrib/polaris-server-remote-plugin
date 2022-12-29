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

package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
)

type rateLimiter struct {
}

func (s *rateLimiter) Ping(_ context.Context, request *api.PingRequest) (*api.PongResponse, error) {
	return &api.PongResponse{}, nil
}

func (r *rateLimiter) Call(_ context.Context, request *api.Request) (*api.Response, error) {
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
		log.Fatal("fail to marshal response data")
		return nil, err
	}

	log.Info("success response", "response", response, "request", request)
	return response, nil
}

func init() {
	logger, err := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()
	if err != nil {
		panic(err)
	}
	log.SetDefaultLoggerWithZap(logger, "plugin-server-v3")
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
