package logger

import (
	"errors"
	"strings"
	"time"
)

/* 自定义日志库
	功能如下：
	- 支持往不同地方写日志
	- 日志分级别
  		- Debug
  		- Info
  		- Warning
  		- Error
  		- Fatal
	- 日志支持开关控制
	- 日志记录包含时间、行号、文件名、日志级别、日志信息
	- 日志文件要切割
*/

// Logevel 自定义日志级别类型
type Logevel uint8

// 定义日志级别
const (
	DEBUG Logevel = iota + 1
	INFO
	WARNING
	ERROR
	FATAL
)

// parseLogLevel
func parseLogLevel(level string) (Logevel, error) {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		return DEBUG, errors.New("未知的日志格式")
	}
}

// 日志级别转字符串
func unparseLogLevel(level Logevel) (string, error) {
	switch level {
	case DEBUG:
		return "DEBUG", nil
	case INFO:
		return "INFO", nil
	case WARNING:
		return "WARNING", nil
	case ERROR:
		return "ERROR", nil
	case FATAL:
		return "FATAL", nil
	default:
		return "DEBUG", errors.New("未知的日志格式")
	}
}

// 格式化时间
func formatTime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}
