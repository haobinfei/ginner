package common

import (
	"fmt"
	"os"
	"time"

	"github.com/haobinfei/ginner/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.SugaredLogger

func InitLogger() {
	nowTime := time.Now()

	infoLogFileName := fmt.Sprintf("%s/Info_%04d-%02d-%02d.log", config.Conf.Logs.Path, nowTime.Year(), nowTime.Month(), nowTime.Day())
	errorLogFileName := fmt.Sprintf("%s/Error_%04d-%02d-%02d.log", config.Conf.Logs.Path, nowTime.Year(), nowTime.Month(), nowTime.Day())

	// 配置日志格式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "time",
		FunctionKey: "func",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2003-01-02 15:04:05"))
		},
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 配置日志级别
	highPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		if config.Conf.Logs.Level >= zap.ErrorLevel {
			return l >= config.Conf.Logs.Level
		}
		return l >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return config.Conf.Logs.Level <= l && l < zap.ErrorLevel
	})

	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFileName,
		MaxSize:    config.Conf.Logs.MaxSize,
		MaxBackups: config.Conf.Logs.MaxBackups,
		MaxAge:     config.Conf.Logs.Maxage,
		LocalTime:  false,
		Compress:   config.Conf.Logs.Compress,
	})
	infoFileCore := zapcore.NewCore(encoder, zap.CombineWriteSyncers(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)

	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFileName,
		MaxSize:    config.Conf.Logs.MaxSize,
		MaxBackups: config.Conf.Logs.MaxBackups,
		MaxAge:     config.Conf.Logs.Maxage,
		LocalTime:  false,
		Compress:   config.Conf.Logs.Compress,
	})
	errorFileCore := zapcore.NewCore(encoder, zap.CombineWriteSyncers(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	logger := zap.New(zapcore.NewTee(infoFileCore, errorFileCore), zap.AddCaller())
	Log = logger.Sugar()
	Log.Info("初始化zap日志完成!")
}
