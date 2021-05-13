package conf

import (
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化 logger

// ConfigDef ConfigDef
type ConfigDef struct {
	FilePath   string `json:"file_path" yaml:"file_path"`     // 日志路径
	MaxSize    int    `json:"max_size" yaml:"max_size"`       // 单个日志最大的文件大小. 单位: MB
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` // 日志文件最多保存多少个备份
	MaxAge     int    `json:"max_age" yaml:"max_age"`         // 文件最多保存多少天
	Compress   bool   `json:"compress" yaml:"compress"`       // 是否压缩备份文件
	Console    bool   `json:"console" yaml:"console"`         // 是否命令行输出，开发环境可以使用
	Level      string `json:"level" yaml:"level"`             // 输出的日志级别, 值：debug,info,warn,error,panic,fatal
}

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	// InfoLevel is the default logging priority.
	"info": zapcore.InfoLevel,
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	"warn": zapcore.WarnLevel,
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	"error": zapcore.ErrorLevel,
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	"dpanic": zapcore.DPanicLevel,
	// PanicLevel logs a message, then panics.
	"panic": zapcore.PanicLevel,
	// FatalLevel logs a message, then calls os.Exit(1).
	"fatal": zapcore.FatalLevel,
}

// Logger Logger
type Logger struct {
	*zap.SugaredLogger
}

// NewLogger NewLogger
func NewLogger(confDef ConfigDef, opts ...zap.Option) *Logger {
	xLogTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}

	hookInfo := lumberjack.Logger{
		Filename:   confDef.FilePath + ".info.log",
		MaxSize:    confDef.MaxSize,
		MaxBackups: confDef.MaxBackups,
		MaxAge:     confDef.MaxAge,
		Compress:   confDef.Compress,
	}

	hookError := lumberjack.Logger{
		Filename:   confDef.FilePath + ".error.log",
		MaxSize:    confDef.MaxSize,
		MaxBackups: confDef.MaxBackups,
		MaxAge:     confDef.MaxAge,
		Compress:   confDef.Compress,
	}

	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		configLevel, ok := levelMap[strings.ToLower(confDef.Level)]
		if !ok {
			return lvl >= zapcore.InfoLevel
		}
		if lvl >= configLevel {
			return true
		}
		return false
	})

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	jsonErr := zapcore.AddSync(&hookError)
	jsonInfo := zapcore.AddSync(&hookInfo)

	// Optimize the xLog output for machine consumption and the console output
	// for human operators.
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "msg"
	encoderConfig.CallerKey = "path"
	encoderConfig.TimeKey = "time"
	//encoderConfig.CallerKey = "path" // 原定的path字段含义太多，建议还是分开，然后log调用的地方就叫caller
	encoderConfig.EncodeTime = xLogTimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	xLogEncoder := zapcore.NewJSONEncoder(encoderConfig)

	var allCore []zapcore.Core
	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the cores together.
	if confDef.FilePath != "" {
		errCore := zapcore.NewCore(xLogEncoder, jsonErr, errPriority)
		infoCore := zapcore.NewCore(xLogEncoder, jsonInfo, infoPriority)
		allCore = append(allCore, errCore, infoCore)
	}

	if confDef.Console {
		consoleDebugging := zapcore.Lock(os.Stdout)
		allCore = append(allCore, zapcore.NewCore(xLogEncoder, consoleDebugging, infoPriority))
	}

	core := zapcore.NewTee(allCore...)
	opts = append(opts, zap.AddCaller())
	logger := zap.New(core).WithOptions(opts...).Sugar()
	defer logger.Sync()
	logger.Infow("test", "test", "test")
	return &Logger{logger}
}
