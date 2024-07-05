package common

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewZaplogger(logLevel string) *zap.Logger {
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
		// TimeKey 指定时间戳在日志条目中的键名
		TimeKey: "time",
		// LevelKey 指定日志级别在日志条目中的键名
		LevelKey: "level",
		// NameKey 指定记录器名称在日志条目中的键名
		NameKey: "logger",
		// CallerKey 指定调用者信息在日志条目中的键名（如文件名、行号、函数名）
		CallerKey: "caller",
		// MessageKey 指定消息内容在日志条目中的键名
		MessageKey: "message",
		// FunctionKey 指定函数名在日志条目中的键名（注意：zap 标准配置中通常不包含此键）
		FunctionKey: "func",
		// StacktraceKey 指定堆栈跟踪信息在日志条目中的键名
		StacktraceKey: "stacktrace",
		// LineEnding 指定日志条目的行结束符，通常为 \n 或 \r\n
		LineEnding: zapcore.DefaultLineEnding,
		// EncodeLevel 指定如何编码日志级别，这里使用小写并带有颜色的编码器
		EncodeLevel: zapcore.LowercaseColorLevelEncoder,
		// EncodeTime 指定如何编码时间戳，这里使用指定的时间格式 "2006-01-02 15:04:05.000"
		EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		// EncodeCaller 指定如何编码调用者信息，这里使用简短的调用者编码器
		EncodeCaller: zapcore.ShortCallerEncoder,
		// EncodeDuration 指定如何编码持续时间，这里使用字符串形式的编码器
		EncodeDuration: zapcore.StringDurationEncoder,
		// EncodeName 指定如何编码记录器名称，这里使用完整的记录器名称编码器
		EncodeName: zapcore.FullNameEncoder,
	}

	// 打印终端日志的core
	zapCoreConsole := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		zapcore.AddSync(os.Stdout),               // 输出到文件
		//zapcore.AddSync(zapcore.Lock(zapcore.AddCaller())), // 输出到文件加调用者信息
		zapcore.DebugLevel, // 日志级别
	)

	// 打印文件日志的core
	core := zapcore.NewTee(
		zapCoreConsole)

	return zap.New(core, zap.AddCaller())
}
