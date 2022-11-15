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
	"fmt"

	"github.com/emicklei/go-restful/v3"
	restfulspec "github.com/polarismesh/go-restful-openapi/v2"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
)

// Resource plugin client resources handler.
type Resource struct {
}

// NewResource returns a new Resource
func NewResource() *Resource {
	return &Resource{}
}

// Plugin plugin model
type Plugin struct {
	Name   string  `json:"name"`
	Config *Config `json:"config"`
}

// WebService return restful web service
func (resource *Resource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/admin/plugins").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	tags := []string{"plugin"}

	ws.Route(
		ws.GET("").To(resource.findAll).
			Doc("get all plugin").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Writes([]Plugin{}).
			Returns(200, "OK", []Plugin{}).
			DefaultReturns("OK", []Plugin{}),
	)

	ws.Route(
		ws.PUT("{name}/disable").To(resource.disable).
			Doc("disable the given name plugin").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Returns(200, "OK", nil).
			DefaultReturns("OK", nil),
	)

	ws.Route(
		ws.PUT("{name}/enable").To(resource.disable).
			Doc("enable the given name plugin").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Returns(200, "OK", nil).
			DefaultReturns("OK", nil),
	)

	return ws
}

// findAll find and return all plugins
//
// GET http://localhost:9050/admin/plugins
func (resource *Resource) findAll(_ *restful.Request, res *restful.Response) {
	var plugins []Plugin
	for name, plugin := range factory.pluginSet {
		plugins = append(plugins, Plugin{
			Name:   name,
			Config: plugin.Config(),
		})
	}
	_ = res.WriteEntity(plugins)
}

// disable disable one plugin by plugin name.
//
// PUT http://localhost:9050/admin/plugins/{name}/disable
func (resource *Resource) disable(req *restful.Request, res *restful.Response) {
	name := req.PathParameter("name")
	if name == "" {
		responseError(res, 400, fmt.Errorf("plugin name is required"))
		return
	}

	plugin := factory.Get(name)
	if plugin == nil {
		responseError(res, 404, fmt.Errorf("plugin with name %s is not exists", name))
		return
	}

	err := plugin.Disable()
	if err != nil {
		responseError(res, 500, fmt.Errorf("fail to disable plugin: %w", err))
		return
	}

	_ = res.WriteEntity(map[string]interface{}{"success": "true"})
	return
}

// enable enable one plugin by plugin name.
//
// PUT http://localhost:9050/admin/plugins/{name}/enable
func (resource *Resource) enable(req *restful.Request, res *restful.Response) {
	name := req.PathParameter("name")
	if name == "" {
		responseError(res, 400, fmt.Errorf("plugin name is required"))
		return
	}

	plugin := factory.Get(name)
	if plugin == nil {
		responseError(res, 404, fmt.Errorf("plugin with name %s is not exists", name))
		return
	}

	err := plugin.Enable()
	if err != nil {
		responseError(res, 500, fmt.Errorf("fail to enable plugin: %w", err))
		return
	}
	_ = res.WriteEntity(map[string]interface{}{"success": "true"})
	return
}

// ResError error response
type ResError struct {
	Err error `json:"err"`
}

func responseError(res *restful.Response, code int, err error) {
	log.Error("admin api returns a error response", "error_response", err)
	_ = res.WriteHeaderAndJson(code, map[string]interface{}{"error": err.Error()}, restful.MIME_JSON)
	return
}
