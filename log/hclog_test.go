package log

import (
	"testing"
)

func Test_hcLoggerWrapper(t *testing.T) {
	Debug("test log", "who", "programmer")
	Info("test log", "who", "programmer")
	Warn("test log", "who", "programmer")
	Error("test log", "who", "programmer")
	Fatal("test log", "who", "programmer")
}
