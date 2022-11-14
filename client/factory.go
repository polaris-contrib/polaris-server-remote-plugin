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

import "sync"

// factory
var factory = NewFactory()

// Factory plugin client factory.
type Factory struct {
	// pluginSet plugin named client set.
	pluginSet map[string]*Client
	// locker locker for pluginSet update/delete
	locker sync.Mutex
}

// NewFactory returns a new Factory.
func NewFactory() *Factory {
	return &Factory{pluginSet: make(map[string]*Client)}
}

// Get return client by name, return nil if the named plugin not exists.
func (f *Factory) Get(name string) *Client {
	client, ex := f.pluginSet[name]
	if !ex {
		return nil
	}
	return client
}

// Register called by plugin client and start up the plugin main process.
func Register(config *Config) (*Client, error) {
	factory.locker.Lock()
	defer factory.locker.Unlock()

	name := config.Name
	if c, ok := factory.pluginSet[name]; ok {
		return c, nil
	}

	c, err := newClient(config)
	if err != nil {
		return nil, err
	}
	factory.pluginSet[name] = c

	return c, nil
}
