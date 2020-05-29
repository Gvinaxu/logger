package main

import (
	logger "github.com/gavlnxu/logger/core"
)

func main() {
	initLogger("file", "./logs", "log_server", "info")
	Run()
}

func initLogger(name, logPath, logName string, level string) (err error) {
	conf := map[string]string{
		"log_path":       logPath,
		"log_name":       logName,
		"log_level":      level,
		"log_split_type": "size",
	}
	err = logger.InitLogger(name, conf)
	if err != nil {
		return
	}

	logger.Info("init logger success")
	return
}

func Run() {
	for {
		logger.Info("log_server is running")
		//time.Sleep(time.Millisecond)
	}
}
