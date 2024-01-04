package bootstrap

import (
	"blog/app"
	"blog/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type DebugLevel struct {
}

func (*DebugLevel) Enabled(level zapcore.Level) bool {
	return level < zap.ErrorLevel
}

func SetupLogger() {
	var logger *zap.Logger

	if config.App.Debug {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = cfg.Build()
	} else {
		debugFile := &lumberjack.Logger{
			Filename:   "./logs/debug.log", // 文件位置
			MaxSize:    50,                 // megabytes，M 为单位，达到这个设置数后就进行日志切割
			MaxBackups: 3,                  // 保留旧文件最大份数
			MaxAge:     31,                 // 旧文件最大保存天数
			Compress:   false,              // 是否压缩日志归档，默认不压缩
		}

		errFile := &lumberjack.Logger{
			Filename:   "./logs/error.log", // 文件位置
			MaxSize:    50,                 // megabytes，M 为单位，达到这个设置数后就进行日志切割
			MaxBackups: 3,                  // 保留旧文件最大份数
			MaxAge:     31,                 // 旧文件最大保存天数
			Compress:   false,              // 是否压缩日志归档，默认不压缩
		}

		defer func(debugFile *lumberjack.Logger, errFile *lumberjack.Logger) {
			_ = debugFile.Close()
			_ = errFile.Close()
		}(debugFile, errFile)

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("[2006/01/02 15:04:05]")
		fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)

		debugLevel := &DebugLevel{}

		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(debugFile), debugLevel),
			zapcore.NewCore(fileEncoder, zapcore.AddSync(errFile), zap.ErrorLevel),
		)

		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	defer logger.Sync()

	app.Logger = logger
}
