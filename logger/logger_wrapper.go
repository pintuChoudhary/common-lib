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

func Fatal(msg string, fields ...any) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Fatal(msg, convertFields(fields...)...)
}

func convertFields(items ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(items))
	for i, v := range items {
		fields = append(fields, zap.Any(fmt.Sprintf("arg_%d", i+1), v))
	}
	return fields
}

func Infof(args ...interface{}) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Info(formatArgs(args...))
}

func Errorf(args ...interface{}) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Error(formatArgs(args...))
}

func Debugf(args ...interface{}) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Debug(formatArgs(args...))
}

func Warnf(args ...interface{}) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Warn(formatArgs(args...))
}

func Fatalf(args ...interface{}) {
	if zapLogger == nil {
		panic("logger not initialized")
	}
	zapLogger.Fatal(formatArgs(args...))
}

func formatArgs(args ...interface{}) string {
	var str string
	for _, arg := range args {
		str += fmt.Sprintf("%+v ", arg)
	}
	return str
}
