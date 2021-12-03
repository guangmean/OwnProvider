package loger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	LOG_LEVEL_ERROR = "ERR"
	LOG_LEVEL_WARN  = "WARN"
	LOG_LEVEL_INFO  = "INFO"
)

func WriteLog(logLevel string, logStr interface{}) {
	logName := time.Now().Format("2006-01-02") + "_ownprovider_inner.log"
	doLog(logLevel, logName, logStr)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func doLog(logLevel string, logName string, logStr interface{}) {
	filePath := "/tmp/commlog/"
	logPath := filePath + logName
	if exist, _ := PathExists(filePath); exist == false {
		fmt.Printf("no dir![%v]\n", filePath)
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("open file failed!")
	} else {
		defer f.Close()
	}
	logger := log.New(f, "DEBUG", log.LstdFlags)
	funcName, _, line, ok := runtime.Caller(1)
	if ok {
		logger.SetPrefix(runtime.FuncForPC(funcName).Name() + "|")
		logger.Printf("|%s|line:%d|%v", logLevel, line, logStr)
	} else {
		logger.Printf("|%s|%v", logLevel, logStr)
	}
}
