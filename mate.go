package zap_mate

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//filaname : config files of zap
//section  : Config options
func NewLogger(filename, section string) *zap.Logger {

	config := viper.New()
	config.SetConfigFile(filename)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	logSec := config.Sub(section)

	return zap.New(zapcore.NewCore(
		func() zapcore.Encoder {
			encoderConfig := zapcore.EncoderConfig{
				MessageKey:     "msg",
				LevelKey:       "level",
				TimeKey:        "time",
				NameKey:        "logger",
				CallerKey:      "file",
				FunctionKey:    "func",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.TimeEncoderOfLayout(logSec.GetString("time-format")),
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}
			if logSec.GetBool("level-color") {
				encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			}
			if logSec.GetString("encoder") == "json" {
				return zapcore.NewJSONEncoder(encoderConfig)
			}
			return zapcore.NewConsoleEncoder(encoderConfig)
		}(),

		zapcore.NewMultiWriteSyncer(
			func() []zapcore.WriteSyncer {
				writes := []zapcore.WriteSyncer{zapcore.AddSync(
					&lumberjack.Logger{
						Filename:   logSec.GetString("file-name"),
						MaxSize:    logSec.GetInt("max-size"),
						MaxBackups: logSec.GetInt("max-count"),
						MaxAge:     logSec.GetInt("max-age"),
						Compress:   false,
						LocalTime:  true,
					})}
				if logSec.GetBool("stdout") {
					writes = append(writes, zapcore.AddSync(os.Stdout))
				}
				return writes
			}()...,
		),
		//DISABLE
		//zapcore.NewMultiWriteSyncer(
		//
		//	func() []zapcore.WriteSyncer {
		//		hook, err := rotatelogs.New(
		//			logSec.GetString("file-name")+".%Y-%m-%d",
		//			rotatelogs.WithLinkName(logSec.GetString("file-name")),
		//			rotatelogs.WithMaxAge(time.Hour*24*time.Duration(logSec.GetInt("max-age"))),
		//			rotatelogs.WithRotationTime(time.Hour*24),
		//			rotatelogs.WithRotationSize(logSec.GetInt64("max-size")*1024*1024),
		//			//rotatelogs.WithRotationCount(logSec.GetUint("max-count")),
		//		)
		//		if err != nil {
		//			panic(hook)
		//		}
		//		writes := []zapcore.WriteSyncer{zapcore.AddSync(hook)}
		//		if logSec.GetBool("stdout") {
		//			writes = append(writes, zapcore.AddSync(os.Stdout))
		//		}
		//		return writes
		//	}()...,
		//),
		zapcore.Level(logSec.GetInt("level")),
	),

		zap.AddStacktrace(zapcore.Level(logSec.GetInt("stack-trace-level"))),

		zap.WithCaller(logSec.GetBool("caller")),
	)

}
