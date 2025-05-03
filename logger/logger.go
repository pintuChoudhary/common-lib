package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	zapLogger *zap.Logger
	once      sync.Once
)

type Config struct {
	ServiceName string // e.g. "auth-service"
	LogToFile   bool
	LogDir      string // optional; if empty, logs only to console
	IsProd      bool
}

// sanitizeFileName replaces dots/colons with underscores
func sanitizeFileName(service string) string {
	timestamp := time.Now().Format("2006_01_02_15_04_05")
	service = strings.ReplaceAll(service, ".", "_")
	service = strings.ReplaceAll(service, ":", "_")
	return fmt.Sprintf("%s_%s.log", service, timestamp)
}

// Init sets up zap logger with optional file logging
func Init(cfg Config) {
	once.Do(func() {
		var core zapcore.Core
		encoderCfg := zap.NewProductionEncoderConfig()
		if !cfg.IsProd {
			encoderCfg = zap.NewDevelopmentEncoderConfig()
		}
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.LevelKey = "L"
		encoderCfg.NameKey = "logger"
		encoderCfg.CallerKey = "C"
		encoderCfg.MessageKey = "M"
		encoderCfg.StacktraceKey = "S"
		encoderCfg.LineEnding = zapcore.DefaultLineEnding
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeDuration = zapcore.SecondsDurationEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

		// Default output to console
		cores := []zapcore.Core{
			zapcore.NewCore(jsonEncoder, zapcore.Lock(os.Stdout), zap.InfoLevel),
		}

		if cfg.LogToFile && cfg.LogDir != "" {
			fileName := sanitizeFileName(cfg.ServiceName)
			fullPath := filepath.Join(cfg.LogDir, fileName)

			writer := zapcore.AddSync(&lumberjack.Logger{
				Filename:   fullPath,
				MaxSize:    10,
				MaxBackups: 1,
				MaxAge:     0,
				Compress:   false,
			})
			fileCore := zapcore.NewCore(jsonEncoder, writer, zap.InfoLevel)
			cores = append(cores, fileCore)
		}

		core = zapcore.NewTee(cores...)
		zapLogger = zap.New(core, zap.AddCaller())
	})
}

func L() *zap.Logger {
	return zapLogger
}
