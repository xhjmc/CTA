package logs_test

import (
	"cta/common/logs"
	"testing"
)

func TestStandardLogger(t *testing.T) {
	logs.Info("this is info", "Info")
	logs.Warn("this is warn", "Warn")
	func() {
		defer func() {
			recover()
		}()
		logs.Panic("this is panic", "Panic")
	}()
	logs.Fatal("this is fatal", "Fatal")
}

func TestStandardLoggerFormat(t *testing.T) {
	logs.Infof("this is info %s", "Infof")
	logs.Warnf("this is warn %s", "Warnf")
	func() {
		defer func() {
			recover()
		}()
		logs.Panicf("this is panic %s", "Panicf")
	}()
	logs.Fatalf("this is fatal %s", "Fatalf")
}
