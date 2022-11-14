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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
)

const (
	// RumModelCompanion companion mode, plugin running by polaris-server
	RumModelCompanion = "companion"
	// RumModelRemote remote mode, plugin running by users
	RumModelRemote = "remote"
)

// RemoteConfig remote plugin config.
type RemoteConfig struct {
	// Address GRPC Service Address
	Address string `yaml:"address" json:"address"`
}

// CompanionConfig companion plugin config.
type CompanionConfig struct {
	// Path is the plugin absolute file path to load.
	Path string `yaml:"path" json:"path"`
	// MaxProcs the max proc number, current plugin can use.
	MaxProcs int `yaml:"max-procs" json:"max_procs"`
	// Args plugin args
	Args []string `yaml:"args" json:"args"`
}

// Config remote plugin config
type Config struct {
	// Name is the plugin unique and exclusive name
	Name string `yaml:"name" json:"name"`
	// Mode is the plugin serverImp running mode, support local and remote.
	Mode string `yaml:"mode" json:"mode"`
	// Remote remote plugin config
	Remote RemoteConfig `yaml:"remote" json:"remote"`
	// Companion companion plugin config
	Companion CompanionConfig `yaml:"companion" json:"companion"`
	// Logger logger instance, use hclog as default
	Logger log.Logger
}

// repairConfig repairs config.
func (c *Config) repairConfig() {
	if c.Companion.MaxProcs == 0 {
		c.Companion.MaxProcs = 1
	}

	if c.Companion.MaxProcs == 0 && c.Companion.MaxProcs >= 4 {
		c.Companion.MaxProcs = 4
	}
}

// pluginLoadPath the path where plugin loading.
func (c *Config) pluginLoadPath() (string, error) {
	fullPath := c.Companion.Path
	if fullPath == "" {
		// Use plugin name and using relative path to load plugin.
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return "", fmt.Errorf("fail to find worksapce: %w", err)
		}
		fullPath = path.Join(dir, c.Name)
	}

	if strings.HasPrefix(fullPath, "..") {
		workspace, _ := os.Getwd()
		fullPath = path.Join(workspace, fullPath)
	}

	if _, err := os.Stat(fullPath); err != nil {
		return "", fmt.Errorf("check plugin file stat error: %w", err)
	}
	return fullPath, nil
}
