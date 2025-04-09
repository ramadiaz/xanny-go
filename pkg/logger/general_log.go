package logger

import (
	"log"
	"os"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

var (
	infoLogger    = log.New(os.Stdout, Blue+"[INFO] "+Reset, log.Ldate|log.Ltime)
	warningLogger = log.New(os.Stdout, Yellow+"[WARNING] "+Reset, log.Ldate|log.Ltime)
	errorLogger   = log.New(os.Stderr, Red+"[ERROR] "+Reset, log.Ldate|log.Ltime)
	panicLogger   = log.New(os.Stderr, Red+"[PANIC] "+Reset, log.Ldate|log.Ltime)
)

func Info(msg string, args ...interface{}) {
	infoLogger.Printf(msg, args...)
}

func Warning(msg string, args ...interface{}) {
	warningLogger.Printf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	errorLogger.Printf(msg, args...)
}

func PanicError(msg string, args ...interface{}) {
	panicLogger.Printf(msg, args...)
	panic("something went wrong, check panic log")
}
