package logs

import (
	"testing"
)

func TestStandardLogger(t *testing.T) {
	Info("this is info", "Info")
	Warn("this is warn", "Warn")
	func() {
		defer func() {
			recover()
		}()
		Panic("this is panic", "Panic")
	}()
	Fatal("this is fatal", "Fatal")
}

func TestStandardLoggerFormat(t *testing.T) {
	Infof("this is info %s", "Infof")
	Warnf("this is warn %s", "Warnf")
	func() {
		defer func() {
			recover()
		}()
		Panicf("this is panic %s", "Panicf")
	}()
	Fatalf("this is fatal %s", "Fatalf")
}