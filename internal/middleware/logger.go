package middleware

import (
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/phonsing-Hub/GoLang/internal/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var LoggerInstance *zap.Logger

func FiberAccessLogger() fiber.Handler {
	return fiberlogger.New(fiberlogger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	})
}

func ZapLogger() fiber.Handler {
	var core zapcore.Core
	var level zapcore.Level

	switch config.Env.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	if config.Env.Development {
		core = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	} else {
		// prod mode: log ลงไฟล์แบบหมุน
		_ = os.MkdirAll("logs", 0755)

		logWriter := &lumberjack.Logger{
			Filename:   "logs/access.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}
		core = zapcore.NewCore(encoder, zapcore.AddSync(logWriter), level)
	}

	LoggerInstance = zap.New(core)

	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Since(start)

		LoggerInstance.Info("request",
			zap.String("ip", c.IP()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", stop),
		)

		return err
	}
}