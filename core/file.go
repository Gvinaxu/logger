package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileLogger struct {
	level         int
	logPath       string
	logName       string
	file          *os.File
	warnFile      *os.File
	logDataChan   chan *LogData
	logSplitType  int
	lastSplitHour int
	logSplitSize  int64 // kb
}

func NewFileLogger(config map[string]string) (logger LogInterface, err error) {
	path, ok := config["log_path"]
	if !ok {
		err = fmt.Errorf("not found 'log_path'")
		return
	}
	name, ok := config["log_name"]
	if !ok {
		err = fmt.Errorf("not found 'log_name'")
		return
	}

	logLevel, ok := config["log_level"]
	if !ok {
		logLevel = "info"
	}
	logChanSize, ok := config["log_chain_size"]
	if !ok {
		logChanSize = "50000"
	}

	logSplitType := LogSplitTypeHour
	logSplitSize := int64(104857600)
	logSplitStr, ok := config["log_split_type"]
	if !ok {
		logSplitStr = "hour"
	} else {
		if logSplitStr == "size" {
			logSplitSizeStr, ok := config["log_split_size"]
			if !ok {
				logSplitSizeStr = "104857600"
			}
			logSplitSize, err = strconv.ParseInt(logSplitSizeStr, 10, 64)
			if err != nil {
				logSplitSize = 104857600
			}
			logSplitType = LogSplitTypeSize
		}
	}

	chanSize, err := strconv.Atoi(logChanSize)
	if err != nil {
		chanSize = 50000
	}

	level := getLogLevel(logLevel)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if err := __ensurePath(absPath); err != nil {
		return nil, err
	}

	logger = &FileLogger{
		level:         level,
		logPath:       absPath,
		logName:       name,
		logDataChan:   make(chan *LogData, chanSize),
		logSplitType:  logSplitType,
		logSplitSize:  logSplitSize,
		lastSplitHour: time.Now().Hour(),
	}
	logger.Init()
	return logger, nil
}

func __ensurePath(p string) error {
	if fi, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(p, os.ModeDir|os.ModePerm); err != nil {
				return err
			}
		}
	} else if !fi.IsDir() {
		return errors.New(fi.Name() + " exists is not directory")
	}
	return nil
}

func (f *FileLogger) Init() {
	fileName := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed, err: %v", fileName, err))
	}
	f.file = file

	wFileName := fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	wFile, err := os.OpenFile(wFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open faile %s failed, err: %v", fileName, err))
	}
	f.warnFile = wFile

	go f.writeLogBackground()
}

func (f *FileLogger) splitFileHour(warnFile bool) {
	now := time.Now()
	hour := now.Hour()
	if hour == f.lastSplitHour {
		return
	}
	f.lastSplitHour = hour
	var backupFileName string
	var fileName string
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%04d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	file.Close()
	os.Rename(fileName, backupFileName)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) splitFileSize(warnFile bool) {
	file := f.file
	if warnFile {
		file = f.warnFile
	}
	info, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := info.Size()
	if fileSize < f.logSplitSize {
		return
	}

	now := time.Now()
	var backupFileName string
	var fileName string
	if warnFile {
		backupFileName = fmt.Sprintf("%s/%s.log.wf_%04d_%02d_%02d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		fileName = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFileName = fmt.Sprintf("%s/%s.log_%04d_%02d_%02d_%02d_%02d_%02d", f.logPath,
			f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		fileName = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}

	file.Close()
	os.Rename(fileName, backupFileName)

	file, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) checkSplitFile(warnFile bool) {
	if f.logSplitType == LogSplitTypeHour {
		f.splitFileHour(warnFile)
		return
	}

	f.splitFileSize(warnFile)

}

func (f *FileLogger) writeLogBackground() {
	for log := range f.logDataChan {
		var file *os.File = f.file
		if log.warnAndFatal {
			file = f.warnFile
		}
		f.checkSplitFile(log.warnAndFatal)
		fmt.Fprintf(file, "%s %s (%s-%s-%d) %s\n", log.timeStr,
			log.levelStr, log.fileName, log.funcName,
			log.lineNum, log.message)
	}
}

func (f *FileLogger) SetLevel(level int) {
	if level <= LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}
	log := NewLogData(LogLevelDebug, format, args...)
	f.logDataChan <- log
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	log := NewLogData(LogLevelTrace, format, args...)

	f.logDataChan <- log
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	log := NewLogData(LogLevelInfo, format, args...)

	f.logDataChan <- log
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	log := NewLogData(LogLevelWarn, format, args...)
	f.logDataChan <- log

}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	log := NewLogData(LogLevelError, format, args...)
	f.logDataChan <- log
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}
	log := NewLogData(LogLevelFatal, format, args...)
	f.logDataChan <- log
}

func (f *FileLogger) Close() {
	f.file.Close()
	f.warnFile.Close()
}
