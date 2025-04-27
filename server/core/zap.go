package core

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"server/config"
	"server/global"
)

// InitLogger 初始化并返回一个基于配置设置的新zap.Logger实例
func InitLogger() *zap.Logger {
	zapConfig := global.Config.Zap

	// 创建一个用于日志输出的writeSyncer
	writeSyncer := getLogWriter(zapConfig)

	// 如果配置了控制台输出，则添加控制台输出
	if zapConfig.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}

	// 创建日志格式化的编码器
	encoder := getEncoder()

	// 根据配置确定日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(zapConfig.Level)); err != nil {
		log.Fatalf("Failed to unmarshal zap level: %v", err)
	}

	// 创建核心和日志实例
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller())

	return logger
}

func getLogWriter(zapConfig config.Zap) zapcore.WriteSyncer {
	// 在打开日志之前，先清空原有内容
	if err := os.Truncate(zapConfig.Filename, 0); err != nil {
		// 如果清空失败，可以选择打印警告，但一般继续运行
		log.Printf("Failed to truncate log file: %v", err)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   zapConfig.Filename,
		MaxSize:    zapConfig.MaxSize,
		MaxBackups: zapConfig.MaxBackups,
		MaxAge:     zapConfig.MaxAge,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
