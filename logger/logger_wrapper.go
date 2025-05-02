package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Info(msg string, fields ...any) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Info(msg, convertFields(fields...)...)
}

func Error(msg string, fields ...any) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Error(msg, convertFields(fields...)...)
}

func Debug(msg string, fields ...any) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Debug(msg, convertFields(fields...)...)
}

func Warn(msg string, fields ...any) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Warn(msg, convertFields(fields...)...)
}

func convertFields(items ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(items))
	for i, v := range items {
		fields = append(fields, zap.Any(fmt.Sprintf("arg_%d", i+1), v))
	}
	return fields
}
