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

package log

import (
	"github.com/hashicorp/go-hclog"
)

// hcLoggerWrapper hclog wrapper.
type hcLoggerWrapper struct {
	hclog.Logger
}

func newHCLoggerWrapper() *hcLoggerWrapper {
	return &hcLoggerWrapper{
		Logger: hclog.New(&hclog.LoggerOptions{
			Level:           hclog.Debug,
			JSONFormat:      true,
			IncludeLocation: true,
		}),
	}
}

// Debug emit a message and key/value pairs at the debug level
func (h *hcLoggerWrapper) Debug(msg string, args ...interface{}) {
	h.Logger.Debug(msg, args)
}

// Info emit a message and key/value pairs at the info level
func (h *hcLoggerWrapper) Info(msg string, args ...interface{}) {
	h.Logger.Info(msg, args)
}

// Warn emit a message and key/value pairs at the warn level
func (h *hcLoggerWrapper) Warn(msg string, args ...interface{}) {
	h.Logger.Warn(msg, args)
}

// Error emit a message and key/value pairs at the error level
func (h *hcLoggerWrapper) Error(msg string, args ...interface{}) {
	h.Logger.Error(msg, args)
}

// Fatal emit a message and key/value pairs at the fatal level
func (h *hcLoggerWrapper) Fatal(msg string, args ...interface{}) {
	h.Error(msg, args...)
}
