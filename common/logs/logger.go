package logs

import (
	"fmt"
	"log"
	"os"
)

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

var logger Logger

func init() {
	SetLogger(&StandardLogger{logger: log.New(os.Stdout, "", log.Lshortfile)})
}

func SetLogger(l Logger) {
	logger = l
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

type StandardLogger struct {
	logger *log.Logger
}

func (l *StandardLogger) Info(args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprint("Info: ", fmt.Sprintln(args...)))
}

func (l *StandardLogger) Warn(args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprint("Warn: ", fmt.Sprintln(args...)))
}

func (l *StandardLogger) Error(args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprint("Error: ", fmt.Sprintln(args...)))
}

func (l *StandardLogger) Fatal(args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprint("Fatal: ", fmt.Sprintln(args...)))
	os.Exit(1)
}

func (l *StandardLogger) Panic(args ...interface{}) {
	s := fmt.Sprint("Panic: ", fmt.Sprintln(args...))
	_ = l.logger.Output(3, s)
	panic(s)
}

func (l *StandardLogger) Infof(format string, args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprintf("Info: "+format, args...))
}

func (l *StandardLogger) Warnf(format string, args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprintf("Warn: "+format, args...))
}

func (l *StandardLogger) Errorf(format string, args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprintf("Error: "+format, args...))
}

func (l *StandardLogger) Fatalf(format string, args ...interface{}) {
	_ = l.logger.Output(3, fmt.Sprintf("Fatal: "+format, args...))
	os.Exit(1)
}

func (l *StandardLogger) Panicf(format string, args ...interface{}) {
	s := fmt.Sprintf("Panic: "+format, args...)
	_ = l.logger.Output(3, s)
	panic(s)
}
