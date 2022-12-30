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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultLogger is the default logger.
var DefaultLogger Logger

func init() {
	logger, _ := zap.Config{
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
	SetDefaultLoggerWithZap(logger, "default")
}

// Logger is the plugin logger interface.
type Logger interface {
	hclog.Logger
	Fatal(msg string, args ...interface{})
}

// Debug omit message and k,v paris at info level.
func Debug(msg string, args ...interface{}) {
	DefaultLogger.Debug(msg, args...)
}

// Info omit message and k,v paris at info level.
func Info(msg string, args ...interface{}) {
	DefaultLogger.Info(msg, args...)
}

// Warn omit message and k,v paris at warn level.
func Warn(msg string, args ...interface{}) {
	DefaultLogger.Warn(msg, args...)
}

// Error omit message and k,v paris at info level.
func Error(msg string, args ...interface{}) {
	DefaultLogger.Error(msg, args...)
}

// Fatal omit message and k,v paris at info level.
func Fatal(msg string, args ...interface{}) {
	DefaultLogger.Fatal(msg, args...)
}
