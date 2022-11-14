package log

import (
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestSetDefaultLoggerWithZap use customization zap logger as default logger.
func TestSetDefaultLoggerWithZap(t *testing.T) {
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

	assert.Nil(t, err)
	stub := gostub.Stub(&DefaultLogger, &ZapLog{logger: logger, name: "test-log"})
	defer stub.Reset()

	Debug("test log", "who", "programmer")
	Info("test log", "who", "programmer")
	Warn("test log", "who", "programmer")
	Error("test log", "who", "programmer")
	Fatal("test log", "who", "programmer")
}
