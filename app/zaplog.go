package app

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

var zapLog *zap.Logger

type Log struct{}

var ZapLog Log

func (log Log) Info(category, msg string) {
	//是否需要重新初始化，生成新的日志文件
	if !needInitLog() {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	lineNum := file + strconv.Itoa(line)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	zapLog.Info(msg,
		zap.String("category", category),
		zap.String("line", lineNum),
		zap.String("log_time", timeStr),
	)
}

func (log Log) Warn(category, msg string) {
	//是否需要重新初始化，生成新的日志文件
	if !needInitLog() {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	lineNum := file + ":" + strconv.Itoa(line)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	zapLog.Warn(msg,
		zap.String("category", category),
		zap.String("line", lineNum),
		zap.String("log_time", timeStr),
	)
}

func (log Log) Error(category, msg string) {
	//是否需要重新初始化，生成新的日志文件
	if !needInitLog() {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	lineNum := file + ":" + strconv.Itoa(line)
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	zapLog.Error(msg,
		zap.String("category", category),
		zap.String("line", lineNum),
		zap.String("log_time", timeStr),
	)
}

func (log Log) Fatal(category, msg string) {
	//是否需要重新初始化，生成新的日志文件
	if !needInitLog() {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	lineNum := file + ":" + strconv.Itoa(line)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	zapLog.Fatal(msg,
		zap.String("category", category),
		zap.String("line", lineNum),
		zap.String("log_time", timeStr),
	)
}

func (log Log) Debug(category, msg string) {
	//是否需要重新初始化，生成新的日志文件
	if !needInitLog() {
		InitLogger()
	}
	_, file, line, _ := runtime.Caller(1)
	lineNum := file + ":" + strconv.Itoa(line)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	zapLog.Debug(msg,
		zap.String("category", category),
		zap.String("line", lineNum),
		zap.String("log_time", timeStr),
	)
}

//设置统一日志路径<elk 会获取此路径日志>
const logpath = "./logs/runtime"

func InitLogger() {
	timeStr := time.Now().Format("2006-01-02")
	fileName := logpath + "/" + timeStr + ".log"
	existBool, _ := isFileExist(fileName)
	if !existBool {
		//创建目录
		err := os.MkdirAll(logpath, os.ModePerm)
		if err != nil {
			log.Panic("日志文件夹创建失败")
		}
		//创建文件
		f, err := os.Create(fileName)
		if err != nil {
			log.Panic("日志文件创建失败")
		}
		f.Close()
	}

	hook := lumberjack.Logger{
		Filename:   fileName, // 日志文件路径
		MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,       // 日志文件最多保存多少个备份
		MaxAge:     7,        // 文件最多保存多少天
		Compress:   true,     // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 设置初始化字段
	//获取服务名称
	filed := zap.Fields(zap.String("serviceName", Config.Server.Name))
	// 构造日志
	logger := zap.New(core, filed)

	zapLog = logger
}

func needInitLog() bool {
	timeStr := time.Now().Format("2006-01-02")
	fileName := logpath + "/" + timeStr + ".log"
	existBool, _ := isFileExist(fileName)
	return existBool
}

//判断文件文件夹是否存在
func isFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}
