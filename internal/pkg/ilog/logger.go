package ilog

import "log"

type Logger struct {
	LoggerName string
}

func (p *Logger) Info(v ...interface{}) {
	log.Println("[Info] ", v)
}

func (p *Logger) Warning(v ...interface{}) {
	log.Println("[Warning] ", v)
}

func (p *Logger) Error(v ...interface{}) {
	log.Println("[Errors] ", v)
}

var l *Logger

func InitLogger(name string) {
	log.SetFlags(log.Ldate | log.Ltime)
	l = &Logger{name}
}

func GetLogger() *Logger {
	return l
}
