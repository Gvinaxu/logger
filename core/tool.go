package core

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogData struct {
	warnAndFatal bool
	message      string
	timeStr      string
	levelStr     string
	fileName     string
	funcName     string
	lineNum      int
}

func getLineInfo() (string, string, int) {
	var (
		fileName, funcName string
		lineNum            int
	)
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
		fileName = file
		lineNum = line
	}
	return fileName, funcName, lineNum
}

func NewLogData(level int, format string, args ...interface{}) *LogData {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05:999")
	levelStr := getLevelText(level)
	fileName, funcName, lineNum := getLineInfo()

	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	msg := fmt.Sprintf(format, args...)
	logData := &LogData{
		message:      msg,
		timeStr:      nowStr,
		levelStr:     levelStr,
		fileName:     fileName,
		funcName:     funcName,
		lineNum:      lineNum,
		warnAndFatal: false,
	}
	switch level {
	case LogLevelError, LogLevelWarn, LogLevelFatal:
		logData.warnAndFatal = true
	}
	return logData

}
