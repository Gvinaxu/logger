package core

import (
	"fmt"
	"os"
)

type ConsoleLogger struct {
	level int
}

func NewConsoleLogger(conf map[string]string) (logger LogInterface, err error) {
	logLevel, ok := conf["log_level"]
	if !ok {
		logLevel = "info"
	}
	level := getLogLevel(logLevel)

	logger = &ConsoleLogger{
		level: level,
	}
	return logger, nil
}

func (c *ConsoleLogger) Init() {
	// todo
}

func (c *ConsoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	c.level = level
}

func (c *ConsoleLogger) Debug(format string, args ...interface{}) {
	if c.level > LogLevelDebug {
		return
	}
	logData := NewLogData(LogLevelDebug, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.timeStr,
		logData.levelStr, logData.fileName, logData.funcName,
		logData.lineNum, logData.message)
}

func (c *ConsoleLogger) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}
	logData := NewLogData(LogLevelTrace, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.timeStr,
		logData.levelStr, logData.fileName, logData.funcName,
		logData.lineNum, logData.message)
}

func (c *ConsoleLogger) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}
	logData := NewLogData(LogLevelInfo, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.timeStr,
		logData.levelStr, logData.fileName, logData.funcName,
		logData.lineNum, logData.message)
}

func (c *ConsoleLogger) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}
	logData := NewLogData(LogLevelWarn, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.timeStr,
		logData.levelStr, logData.fileName, logData.funcName,
		logData.lineNum, logData.message)
}

func (c *ConsoleLogger) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}
	logData := NewLogData(LogLevelError, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.timeStr,
		logData.levelStr, logData.fileName, logData.funcName,
		logData.lineNum, logData.message)
}

func (c *ConsoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}
	logData := NewLogData(LogLevelFatal, format, args...)
	fmt.Fprintf(os.Stdout, "%s %s (%s:%s:%d) %s\n", logData.timeStr,
		logData.levelStr, logData.fileName, logData.funcName,
		logData.lineNum, logData.message)
}

func (c *ConsoleLogger) Close() {
	// todo
}
