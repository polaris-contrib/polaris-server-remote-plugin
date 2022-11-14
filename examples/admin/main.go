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
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common/client"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
)

func main() {
	if _, err := client.Register(
		&client.Config{
			Name: "remote-rate-limit-server-v1",
			Mode: client.RumModelCompanion,
			Local: client.CompanionConfig{
				MaxProcs: 1,
				Path:     "./remote-rate-limit-server-v1",
			},
		},
	); err != nil {
		log.Fatal("server-v1 register failed", "error", err.Error())
		return
	}

	if _, err := client.Register(
		&client.Config{
			Name: "remote-rate-limit-server-v2",
			Mode: client.RumModelCompanion,
			Local: client.CompanionConfig{
				Path: "./remote-rate-limit-server-v2",
			},
		},
	); err != nil {
		log.Fatal("server-v2 register failed", "error", err.Error())
	}

	restful.DefaultContainer.Add(client.NewResource().WebService())
	adminPort := 9050

	log.Info(fmt.Sprintf("request the admin api using http://localhost:%d", adminPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", adminPort), nil); err != nil {
		log.Fatal("plugin admin serve error", "error", err)
	}
}
