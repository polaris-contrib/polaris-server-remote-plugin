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

	pluginsdk "github.com/polaris-contrib/polaris-server-remote-plugin-common"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/api"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/server"
)

type rateLimiter struct{}

func (s *rateLimiter) Call(_ context.Context, request *api.Request) (*api.Response, error) {
	var rateLimitRequest api.RateLimitPluginRequest
	if err := pluginsdk.UnmarshalRequest(request, &rateLimitRequest); err != nil {
		log.Error("fail to marshal response data", "error", err)
		return nil, err
	}
	response, err := pluginsdk.MarshalResponse(&api.RateLimitPluginResponse{Allow: true})
	if err != nil {
		log.Error("fail to marshal response data", "error", err)
		return nil, err
	}
	log.Info("success response", "response", response, "request", request)
	return response, nil
}

func main() {
	server.Serve(&rateLimiter{})
}
