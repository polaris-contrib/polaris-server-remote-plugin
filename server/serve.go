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

package server

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	pluginsdk "github.com/polaris-contrib/polaris-server-remote-plugin-common"
	"github.com/polaris-contrib/polaris-server-remote-plugin-common/plugins"
)

var (
	// parentPID polaris server main process id.
	parentPID int

	// LivenessCheckInterval the interval of liveness
	LivenessCheckInterval = time.Second

	// logger
	logger hclog.Logger
)

func init() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})
	parentPID = os.Getppid()
}

// Serve is a function used to serve a plugin.
//
// Keep in mind that do not write log messages to console, because go-plugin use stdout、stderr pipe
// as the communication bridge for plugin-client and plugin-server, those log messages may prevent
// the plugin-system from working properly.
//
// The plugin.Serve already contains gRPC Graceful exit logic, so don't implement it additionally.
func Serve(svc plugins.Service) {
	// polaris server will set it.
	pluginName := os.Getenv("PLUGIN_NAME")
	if pluginName == "" {
		pluginName = filepath.Base(os.Args[0])
	}

	// polaris server will set it.
	p := os.Getenv("PLUGIN_PROCS")
	if procs, err := strconv.Atoi(p); err == nil {
		runtime.GOMAXPROCS(procs)
	}

	// logger.Info("plugin ", pluginName, " use cpu numbers: ", runtime.GOMAXPROCS(0))

	go func() {
		checkMainProcessLivenessTimer()
	}()

	// logger.Info("plugin ", pluginName, " is running ...")
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pluginsdk.Handshake,
		Plugins: map[string]plugin.Plugin{
			pluginName: plugins.NewPlugin(svc),
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

// checkMainProcessLivenessTimer check main process still alive for every interval.
//
// Because we do not use the re-attach capability of the go-plugin, once the main process of
// the Polaris service exits, we need to exit the plug-in to prevent many zombie plugins
// from entering the field
func checkMainProcessLivenessTimer() {
	ticker := time.NewTicker(LivenessCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			canIExit()
		}
	}
}

// canIExit check main process is alive or not, exit when main process was exited.
func canIExit() {
	// Process ID 1 is id of the init process, when equals to 1 means current process is zombie process.
	if parentPID == 1 || os.Getppid() != parentPID {
		logger.Error("polaris server has gone, %s will exit ...", filepath.Base(os.Args[0]))
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}

	if _, err := os.FindProcess(parentPID); err != nil {
		logger.Error("can not find polaris server, plugin %s will exit ...", filepath.Base(os.Args[0]))
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
}
