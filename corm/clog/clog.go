package clog

import (
	"log"
	"os"
	"runtime"
	"strconv"
)

//颜色ansi代码
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

//规定打印日志等级
const (
	DebugLevel = iota + 1
	InfoLevel
	ErrorLevel
)

//规定打印的日志语句
var (
	level       = DebugLevel
	sqlLogger   = log.New(os.Stdout, Yellow+"[SQL]"+Reset, 0)
	debugLogger = log.New(os.Stdout, Blue+"[Debug] "+Reset, log.LstdFlags)
	infoLogger  = log.New(os.Stdout, Green+"[Info] "+Reset, log.LstdFlags)
	errorLogger = log.New(os.Stdout, Red+"[Error] "+Reset, log.LstdFlags)
)

// SetLevel 自定义打印日志的等级，默认为info等级（开发的默认是debug等级）
func SetLevel(l int) {
	level = l
}

//Sql 打印Sql语句
func Sql(msg string, args ...any) {
	if len(args) == 0 {
		sqlLogger.Printf(msg + "\n")
	} else {
		sqlLogger.Printf(msg+"%v"+"\n", append([]any{}, args...))
	}
}

//Debug 打印开发者日志
func Debug(msg string, args ...any) {
	if level <= DebugLevel {
		if len(args) == 0 {

		} else {
			debugLogger.Printf(getFileWithLine()+msg+"\n", append([]any{}, args...))
		}
	}
}

//Info 打印Info消息
func Info(msg string, args ...any) {
	if level <= InfoLevel {
		if len(args) == 0 {
			infoLogger.Printf(msg + "\n")
		} else {
			infoLogger.Printf(msg+"\n", append([]any{}, args...))
		}
	}
}

//Error 打印错误信息
func Error(msg string, args ...any) {
	if level <= ErrorLevel {
		if len(args) == 0 {
			errorLogger.Printf(getFileWithLine() + msg + "\n")
		} else {
			errorLogger.Printf(getFileWithLine()+msg+"\n", args...)
		}
	}
}

func getFileWithLine() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return file + ":" + strconv.Itoa(line) + " "
	}
	return ""
}
