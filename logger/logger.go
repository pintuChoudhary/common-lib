package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	zapLogger *zap.Logger
	LogLevel  atomic.Int32
	initOnce  sync.Once
)

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

type Config struct {
	ServiceName string
	LogToFile   bool
	LogDir      string
	IsProd      bool
	LogLevel    zapcore.Level
}

func sanitizeFileName(service string) string {
	timestamp := time.Now().Format("2006_01_02_15_04_05")
	service = strings.ReplaceAll(service, ".", "_")
	service = strings.ReplaceAll(service, ":", "_")
	return fmt.Sprintf("%s_%s.log", service, timestamp)
}

func Init(cfg Config) {
	initOnce.Do(func() {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "time"
		encoderCfg.LevelKey = "level"
		encoderCfg.MessageKey = "msg"
		encoderCfg.CallerKey = ""     // Disable caller field
		encoderCfg.StacktraceKey = "" // Disable stacktrace
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

		jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

		cores := []zapcore.Core{
			zapcore.NewCore(
				jsonEncoder,
				zapcore.Lock(os.Stdout),
				zap.InfoLevel,
			),
		}

		if cfg.LogToFile && cfg.LogDir != "" {
			fullPath := filepath.Join(cfg.LogDir, sanitizeFileName(cfg.ServiceName))

			fileWriter := zapcore.AddSync(&lumberjack.Logger{
				Filename:   fullPath,
				MaxSize:    10,   // MB
				MaxBackups: 5,    // Number of backups
				MaxAge:     30,   // Days
				Compress:   true, // Compress rotated files
			})

			cores = append(cores, zapcore.NewCore(
				jsonEncoder,
				fileWriter,
				zap.InfoLevel,
			))
		}

		core := zapcore.NewTee(cores...)
		zapLogger = zap.New(core)
		SetLevel(cfg.LogLevel)
	})
}

func L() *zap.Logger {
	return zapLogger
}

func SetLevel(level zapcore.Level) {
	LogLevel.Store(int32(level))
}

// GetLevel returns current log level
func GetLevel() zapcore.Level {
	return zapcore.Level(LogLevel.Load())
}

// Helper functions with structured logging
func Info(msg string, fields ...zap.Field) {
	if shouldLog(zapcore.InfoLevel) {
		zapLogger.Info(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if shouldLog(zapcore.ErrorLevel) {
		zapLogger.Error(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if shouldLog(zapcore.DebugLevel) {
		zapLogger.Debug(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if shouldLog(zapcore.WarnLevel) {
		zapLogger.Warn(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if shouldLog(zapcore.FatalLevel) {
		zapLogger.Fatal(msg, fields...)
	}
}

// shouldLog checks if we should log at given level
func shouldLog(level zapcore.Level) bool {
	return level >= zapcore.Level(LogLevel.Load())
}

// Structured field helpers
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func ErrorWrap(err error) zap.Field {
	return zap.Error(err)
}
