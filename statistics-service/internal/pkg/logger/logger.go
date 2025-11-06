package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level      string `mapstructure:"level"`       // 日志级别
	FilePath   string `mapstructure:"file_path"`   // 文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 文件最大大小(MB)
	MaxBackups int    `mapstructure:"max_backups"` // 最大备份数
	MaxAge     int    `mapstructure:"max_age"`     // 保存天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

// InitLogger 初始化日志
func InitLogger(cfg *Config) error {
	// 设置日志级别
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}

	// 编码配置
	enconderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志输出
	cores := []zapcore.Core{}

	// 文件输出
	if cfg.FilePath != "" {
		fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(enconderConfig),
			fileWriteSyncer,
			level,
		)
		cores = append(cores, fileCore)
	}

	// 控制台输出（开发环境）
	consoleEncoder := zapcore.NewConsoleEncoder(enconderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
	cores = append(cores, consoleCore)

	// 创建 Logger
	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	SugarLogger = Logger.Sugar()
	return nil
}

func Sync() {
	_ = SugarLogger.Sync()
}
