package logger_utils

import (
	"go.uber.org/zap"
)

func ConfigureDefaultLogger() {
	logger := zap.NewExample()

	std := zap.NewStdLog(logger)
	std.Print("standard logger wrapper")
}
