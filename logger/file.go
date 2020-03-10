package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// FileLoger 文件输出日志
type FileLoger struct {
	level         Logevel
	filePath      string
	fileName      string
	fileObj       *os.File
	errFileObj    *os.File
	maxFileSize   int64
	splitLogTime  string
	splitLogOnOff bool
}

// FileHandler 文件输出
func FileHandler(level, filePath, fileName string, maxSize int64, splitLogOnOff bool) *FileLoger {
	// 对日志级别进行格式化处理
	loglevel, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	f1 := &FileLoger{
		level:         loglevel,
		filePath:      filePath,
		fileName:      fileName,
		maxFileSize:   maxSize,
		splitLogTime:  time.Now().Format("2006-02-01"),
		splitLogOnOff: splitLogOnOff,
	}
	// 对文件进行打开初始化操作
	err = f1.initFile()
	if err != nil {
		panic(err)
	}
	return f1
}

// 对文件进行初始化操作
func (f *FileLoger) initFile() error {
	fullFilePath := path.Join(f.filePath, f.fileName)
	// 打开文件
	fileObj, err := os.OpenFile(fullFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("日志文件打开失败，err:", err)
		return err
	}
	// 打开错误日志
	errFileObj, err := os.OpenFile(fullFilePath+".error", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("错误日志文件打开失败，err:", err)
		return err
	}
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}

// 比较日志级别
func (f *FileLoger) logEnable(l Logevel) bool {
	return f.level <= l
}

// 切割日志操作
func (f *FileLoger) splitLog(file *os.File, nowStr string) (*os.File, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("get file info failed. err:", err)
		return nil, err
	}
	oldLogName := path.Join(f.filePath, fileInfo.Name())
	newLogName := fmt.Sprintf("%s.backup.%s", oldLogName, nowStr)
	file.Close()
	os.Rename(oldLogName, newLogName)
	newFileObj, err := os.OpenFile(oldLogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return newFileObj, err
}

// 日志切割（按大小）
func (f *FileLoger) splitLogBySize(file *os.File) (*os.File, error) {
	/*
		1、获取现有日志大小和日志名
		2、判断日志大小是否大于设定
		3、对原有日志进行备份,按时间戳备份
		4、关闭原有日志
		5、新打开一个日志
		6、将新打开日志赋值给原有日志变量
	*/
	// 获取时间戳
	nowStr := time.Now().Format("2006-01-02-15040505")
	fileInfo, err := file.Stat()
	fmt.Println(fileInfo)
	if err != nil {
		fmt.Println("get file info failed. err:", err)
		return nil, err
	}
	if fileInfo.Size() >= f.maxFileSize {
		newFileObj, err := f.splitLog(file, nowStr)
		if err != nil {
			fmt.Println("open new file failed,bysize. err:", err)
			return nil, err
		}
		return newFileObj, nil
	}
	return nil, nil
}

// 日志切割（按时间）
func (f *FileLoger) splitLogByTime(file *os.File) (*os.File, error) {
	/*
		1、获取当前的时间
		2、与上一次切割的时间做比较
		3、切割后记录本次切割的时间
	*/
	// 获取当前时间
	nowTime := time.Now()
	nowTimeStr := nowTime.Format("2006-02-01")
	// 设置时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 对上次切割的时间进行转换
	lastSplitTime, _ := time.ParseInLocation("2006-01-02", f.splitLogTime, loc)
	// 对两次时间进行比较
	if nowTime.Sub(lastSplitTime).Hours()/24 < 1 {
		newFileObj, err := f.splitLog(file, nowTimeStr)
		if err != nil {
			fmt.Println("open new file failed,bytime. err:", err)
			return nil, err
		}
		f.splitLogTime = nowTimeStr
		return newFileObj, nil
	}
	return nil, nil
}

// 判断是按时间切割还是按日期切割
func (f *FileLoger) splitByTimeOrLog(file *os.File) *os.File {
	/*
		1、如果splitLogOnOff为true则按大小切割
		2、如果splitLogOnOff为false则按时间切割
	*/
	if f.splitLogOnOff {
		newFileObj, err := f.splitLogBySize(file)
		if err != nil {
			fmt.Println("file split failed. err:", err)
			return nil
		}
		if newFileObj != nil {
			return newFileObj
		}
	} else {
		newFileObj, err := f.splitLogByTime(file)
		if err != nil {
			fmt.Println("file split failed. err:", err)
			return nil
		}
		if newFileObj != nil {
			return newFileObj
		}
	}
	return nil
}

// 日志内容 [时间] [日志级别] [文件名：函数名：行号] 日志内容
func (f *FileLoger) logInfo(loglevel Logevel, format string, a ...interface{}) {
	if f.logEnable(loglevel) {
		msg := fmt.Sprintf(format, a...)
		// 日志级别转字符串
		levelStr, err := unparseLogLevel(loglevel)
		if err != nil {
			fmt.Println(err)
		}
		funcName, fileName, lineNo := getFileInfo(3)

		newFileObj := f.splitByTimeOrLog(f.fileObj)
		if newFileObj != nil {
			f.fileObj = newFileObj
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s-%s-%d] %s\n", formatTime(), levelStr, fileName, funcName, lineNo, msg)
		// defer f.fileObj.Close()
		if loglevel >= ERROR {
			newErrFileObj := f.splitByTimeOrLog(f.errFileObj)
			if newErrFileObj != nil {
				f.errFileObj = newErrFileObj
			}
			fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s-%s-%d] %s\n", formatTime(), levelStr, fileName, funcName, lineNo, msg)
			// defer f.errFileObj.Close()
		}
	}
}

// Debug 定义Debug日志方法
func (f *FileLoger) Debug(format string, a ...interface{}) {
	f.logInfo(DEBUG, format, a...)
}

// Info 定义Info日志方法
func (f *FileLoger) Info(format string, a ...interface{}) {
	f.logInfo(INFO, format, a...)
}

// Warning 定义Warning日志方法
func (f *FileLoger) Warning(format string, a ...interface{}) {
	f.logInfo(WARNING, format, a...)
}

// Error 定义Error日志方法
func (f *FileLoger) Error(format string, a ...interface{}) {
	f.logInfo(ERROR, format, a...)
}

// Fatal 定义Fatal日志方法
func (f *FileLoger) Fatal(format string, a ...interface{}) {
	f.logInfo(FATAL, format, a...)
}
