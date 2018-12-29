package logs

import (
	"fmt"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func GetBlogger() Logger {
	return DefaultLogger
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, args ...interface{})
	Error(v ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, args ...interface{})
	Info(v ...interface{})
	Infof(format string, args ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, args ...interface{})
}

var DefaultLogger = new(StdoutLogger)

type StdoutLogger struct{}

func (logger *StdoutLogger) Debug(v ...interface{}) {
	logger.PrintWithLevel("DEBUG", v...)
}
func (logger *StdoutLogger) Debugf(format string, args ...interface{}) {
	logger.PrintfWithLevel("DEBUG", format, args...)
}
func (logger *StdoutLogger) Error(v ...interface{}) {
	logger.PrintWithLevel("ERROR", v...)
}
func (logger *StdoutLogger) Errorf(format string, args ...interface{}) {
	logger.PrintfWithLevel("ERROR", format, args...)
}
func (logger *StdoutLogger) Fatal(v ...interface{}) {
	logger.PrintWithLevel("FATAL", v...)
}
func (logger *StdoutLogger) Fatalf(format string, args ...interface{}) {
	logger.PrintfWithLevel("FATAL", format, args...)
}
func (logger *StdoutLogger) Info(v ...interface{}) {
	logger.PrintWithLevel("INFO", v...)
}
func (logger *StdoutLogger) Infof(format string, args ...interface{}) {
	logger.PrintfWithLevel("INFO", format, args...)
}
func (logger *StdoutLogger) Warn(v ...interface{}) {
	logger.PrintWithLevel("WARN", v...)
}
func (logger *StdoutLogger) Warnf(format string, args ...interface{}) {
	logger.PrintfWithLevel("WARN", format, args...)
}
func (logger *StdoutLogger) Print(v ...interface{}) {
	log.Output(3, fmt.Sprint(v...))
}
func (logger *StdoutLogger) PrintWithLevel(l string, v ...interface{}) {
	log.Output(3, fmt.Sprint(v...))
}
func (logger *StdoutLogger) PrintfWithLevel(l, format string, args ...interface{}) {
	log.Output(3, fmt.Sprintf(l+" "+format, args...))
}
func (logger *StdoutLogger) Println(v ...interface{}) {
	log.Output(3, fmt.Sprint(v...))
}
func (logger *StdoutLogger) Printf(format string, args ...interface{}) {
	log.Output(3, fmt.Sprintf(format, args...))
}
