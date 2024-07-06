package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewZaplogger(logLevel, LogPath, ErrLogPath string) *zap.Logger {
	// 日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch logLevel {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	// 设置默认字段
	// encoderConfig 是一个 zapcore.EncoderConfig 结构体实例，用于配置日志条目的编码方式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                                                 // 时间字段
		LevelKey:       "level",                                                // 日志级别字段
		NameKey:        "name",                                                 // 日志记录器名称字段
		CallerKey:      "line",                                                 // 调用者字段（代码行号）
		MessageKey:     "msg",                                                  // 日志消息字段
		FunctionKey:    "func",                                                 // 函数名字段
		StacktraceKey:  "stacktrace",                                           // 堆栈跟踪字段
		LineEnding:     zapcore.DefaultLineEnding,                              // 行结束符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                          // 日志级别编码（小写）
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"), // 时间编码格式
		EncodeDuration: zapcore.SecondsDurationEncoder,                         // 持续时间编码
		EncodeCaller:   zapcore.FullCallerEncoder,                              // 调用者编码（全路径）
		EncodeName:     zapcore.FullNameEncoder,                                // 日志记录器名称编码
	}

	// 标准日志轮转
	writer := &lumberjack.Logger{
		Filename:   LogPath, // 日志文件路径
		MaxSize:    1,       // 每个日志文件的最大大小（单位：M）
		MaxBackups: 30,      // 日志文件最多保存多少个备份
		MaxAge:     7,       // 文件最多保存多少天
		LocalTime:  true,    // 开启本地时间
		Compress:   true,    // 开启压缩
	}

	// 错误日志轮转
	errorWriter := &lumberjack.Logger{
		Filename:   ErrLogPath, // 日志文件路径
		MaxSize:    100,        // 每个日志文件的最大大小（单位：M）
		MaxBackups: 30,         // 日志文件最多保存多少个备份
		MaxAge:     7,          // 文件最多保存多少天
		LocalTime:  true,       // 开启本地时间
		Compress:   true,       // 开启压缩
	}

	// 写入日志文件和终端的zapCore.Core
	zapCoreError := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		zapcore.AddSync(errorWriter),          // 输出到文件
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel // 只允许ERROR级别以上的日志输出到文件
		}), // 日志级别
	)

	// 打印终端日志的core
	zapCoreConsole := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		zapcore.AddSync(os.Stdout),               // 输出到文件
		//zapcore.AddSync(zapcore.Lock(zapcore.AddCaller())), // 输出到文件加调用者信息
		zapcore.DebugLevel, // 日志级别
	)

	// 写入日志文件用的CoreFile
	zapCoreFile := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		zapcore.AddSync(writer),               // 输出到文件
		atomicLevel,                           // 日志级别

	)

	//  合并coreFile和console
	core := zapcore.NewTee(
		zapCoreFile,    // 文件输出Core
		zapCoreConsole, // 终端输出终端控制台
		zapCoreError,   // 文件输出错误日志
	)

	return zap.New(core, zap.AddCaller())

}
