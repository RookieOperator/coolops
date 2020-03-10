package logger

import (
	"fmt"
	"path"
	"runtime"
)

// ConsolLogger 定义日志结构体
type ConsolLogger struct {
	loglevel Logevel
}

// SetLevel 设置日志级别
// func (l Logger) SetLevel() {
// }

// StreamHandler 屏幕输出
func StreamHandler(level string) ConsolLogger {
	// 对日志级别进行格式化处理
	loglevel, err := parseLogLevel(level)
	if err != nil {
		fmt.Println(err)
	}
	return ConsolLogger{
		loglevel: loglevel,
	}
}

// 比较日志级别
func (c ConsolLogger) comLogLevel(l Logevel) bool {
	return c.loglevel <= l
}

// 获取文件名，函数名，行号
func getFileInfo(skip int) (fileName, funcName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("获取文件信息失败")
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	return
}

// 日志内容 [时间] [日志级别] [文件名：函数名：行号] 日志内容
func (c ConsolLogger) logInfo(loglevel Logevel, format string, a ...interface{}) {
	if c.comLogLevel(loglevel) {
		msg := fmt.Sprintf(format, a...)
		// 日志级别转字符串
		levelStr, err := unparseLogLevel(loglevel)
		if err != nil {
			fmt.Println(err)
		}
		funcName, fileName, lineNo := getFileInfo(3)
		fmt.Printf("[%s] [%s] [%s-%s-%d] %s\n", formatTime(), levelStr, fileName, funcName, lineNo, msg)
	}
}

// Debug 定义Debug日志方法
func (c ConsolLogger) Debug(format string, a ...interface{}) {
	c.logInfo(DEBUG, format, a...)
}

// Info 定义Info日志方法
func (c ConsolLogger) Info(format string, a ...interface{}) {
	c.logInfo(INFO, format, a...)
}

// Warning 定义Warning日志方法
func (c ConsolLogger) Warning(format string, a ...interface{}) {
	c.logInfo(WARNING, format, a...)
}

// Error 定义Error日志方法
func (c ConsolLogger) Error(format string, a ...interface{}) {
	c.logInfo(ERROR, format, a...)
}

// Fatal 定义Fatal日志方法
func (c ConsolLogger) Fatal(format string, a ...interface{}) {
	c.logInfo(FATAL, format, a...)
}
