// zapLog doc

package zapLog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var logger *zap.Logger

var logConfigFileName = "log.yaml"

func init()  {
	initLog()
}

//initLog doc
//@Description: log初始化
//@Author niejian
//@Date 2021-04-26 14:07:35
func initLog() {

	// 日志配置文件
	projectPath, _ := os.Getwd()
	logConfigPath := fmt.Sprintf("%s%s%s%s%s", projectPath, string(filepath.Separator), "resources",
		string(filepath.Separator), logConfigFileName)

	// 读取配置文件
	//logConfig := parser.LogConfigParser()
	logConfig := watchAndReloadConfig(logConfigPath)

	if nil == logConfig {
		panic("初始化日志配置失败，请检查")
	}

	dir, _ := os.Getwd()
	////dir := getCurrentPath()
	//fmt.Println("路径：" + dir)
	var logPath = dir + string(filepath.Separator) + logConfig.LogPath
	// 判断文件夹是否存在
	if _, err := os.Stat(logPath); err != nil {
		fmt.Println("日志文件夹:" + logPath + ", 不存在，创建")
		// 文件夹不存在，创建
		os.Mkdir(logPath, os.ModePerm)
	}
	logFile := logPath + string(filepath.Separator) + logConfig.LogName
	fmt.Println("日志路径：" + logFile)
	hook := lumberjack.Logger{
		Filename:   logFile, // 日志文件路径
		MaxSize:    logConfig.MaxSize,                                        // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: logConfig.MaxBackups,                                     // 日志文件最多保存多少个备份
		MaxAge:     logConfig.MaxAge,                                         // 文件最多保存多少天
		Compress:   true,                                                     // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	//atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	//config := zap.Config{
	//	Level:            atomicLevel,                                                // 日志级别
	//	Development:      true,                                                // 开发模式，堆栈跟踪
	//	Encoding:         "console",                                           // 输出格式 console 或 json
	//	EncoderConfig:    encoderConfig,                                         // 编码器配置
	//	//InitialFields:    map[string]interface{}{"serviceName": "k8sApi"},     // 初始化字段，如：添加一个服务器名称
	//	OutputPaths:      []string{"stdout", LOGGER_PATH},       // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
	//	ErrorOutputPaths: []string{"stderr"},
	//
	//}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(logConfig.LogLevel))

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	//filed := zap.Fields(zap.String("", ""))
	// 构造日志
	initLogger := zap.New(core, caller, development)

	// 构建日志
	initLogger.Info("log 初始化成功 ")
	logger = initLogger
}

func LOGGER() *zap.Logger {
	return logger
}

func reloadLogger()  {
	initLog()
}
