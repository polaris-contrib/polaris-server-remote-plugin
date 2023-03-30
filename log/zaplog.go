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
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/go-hclog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultZapLogName = "zap"

// ZapLog zap base Logger implementation.
type ZapLog struct {
	logger *zap.Logger
	name   string
}

// SetDefaultLoggerWithZap create a zap based log with the given core.
func SetDefaultLoggerWithZap(logger *zap.Logger, name string) {
	DefaultLogger = &ZapLog{logger: logger, name: name}
}

// Log Emit a message and key/value pairs at a provided log level
func (h *ZapLog) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.NoLevel:
		return
	case hclog.Trace:
		h.Trace(msg, args)
	case hclog.Debug:
		h.Debug(msg, args)
	case hclog.Info:
		h.Info(msg, args)
	case hclog.Warn:
		h.Warn(msg, args)
	case hclog.Error:
		h.Error(msg, args)
	case hclog.Off:
		return
	}
}

// Trace emit a message and key/value pairs at the TRACE level
func (h *ZapLog) Trace(msg string, args ...interface{}) {
	h.logger.Debug(msg, anythingsToZapFields(args...)...)
}

// Debug emit a message and key/value pairs at the DEBUG level
func (h *ZapLog) Debug(msg string, args ...interface{}) {
	h.logger.Debug(msg, anythingsToZapFields(args...)...)
}

// Info emit a message and key/value pairs at the INFO level
func (h *ZapLog) Info(msg string, args ...interface{}) {
	h.logger.Info(msg, anythingsToZapFields(args...)...)
}

// Warn emit a message and key/value pairs at the WARN level
func (h *ZapLog) Warn(msg string, args ...interface{}) {
	h.logger.Warn(msg, anythingsToZapFields(args...)...)
}

// Error emit a message and key/value pairs at the ERROR level
func (h *ZapLog) Error(msg string, args ...interface{}) {
	h.logger.Error(msg, anythingsToZapFields(args...)...)
}

// Fatal emit a message and key/value pairs at the Fatal level
func (h *ZapLog) Fatal(msg string, args ...interface{}) {
	h.logger.Fatal(msg, anythingsToZapFields(args...)...)
}

// IsTrace indicate if TRACE logs would be emitted.
func (h *ZapLog) IsTrace() bool {
	return true
}

// IsDebug indicate if DEBUG logs would be emitted.
func (h *ZapLog) IsDebug() bool {
	return true
}

// IsInfo indicate if INFO logs would be emitted
func (h *ZapLog) IsInfo() bool {
	return true
}

// IsWarn indicate if WARN logs would be emitted.
func (h *ZapLog) IsWarn() bool {
	return true
}

// IsError indicate if ERROR logs would be emitted.
func (h *ZapLog) IsError() bool {
	return true
}

// ImpliedArgs returns With key/value pairs
func (h *ZapLog) ImpliedArgs() []interface{} {
	return []interface{}{}
}

// With creates a sub logger that will always have the given key/value pairs
func (h *ZapLog) With(args ...interface{}) hclog.Logger {
	return h.With(args)
}

// Name returns the Name of the logger
func (h *ZapLog) Name() string {
	return h.name
}

// Named create a logger that will prepend the name string on the front of all messages.
func (h *ZapLog) Named(name string) hclog.Logger {
	newLogger := &ZapLog{name: name}
	return newLogger
}

// ResetNamed reset log name
func (h *ZapLog) ResetNamed(name string) hclog.Logger {
	h.name = defaultZapLogName
	return h
}

// SetLevel updates the level.
func (h *ZapLog) SetLevel(level hclog.Level) {
}

// StandardLogger return a value that conforms to the stdlib log.Logger interface
func (h *ZapLog) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return log.New(h.StandardWriter(opts), "", log.LstdFlags)
}

// StandardWriter Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (h *ZapLog) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return hclog.DefaultOutput
}

func anythingsToZapFields(args ...interface{}) []zap.Field {
	var fields []zapcore.Field
	for i := len(args); i > 0; i -= 2 {
		left := i - 2
		if left < 0 {
			left = 0
		}

		items := args[left:i]

		switch l := len(items); l {
		case 2:
			k, ok := items[0].(string)
			if ok {
				fields = append(fields, zap.Any(k, items[1]))
			} else {
				fields = append(fields, zap.Any(fmt.Sprintf("field-%d", i-1), items[1]))
				fields = append(fields, zap.Any(fmt.Sprintf("field-%d", left), items[0]))
			}
		case 1:
			fields = append(fields, zap.Any(fmt.Sprintf("field-%d", left), items[0]))
		}
	}

	return fields
}
