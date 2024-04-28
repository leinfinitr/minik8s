package log

import (
	"testing"
)

func TestInfoLog(t *testing.T) {
	InfoLog("test", "test")
}

func TestDebugLog(t *testing.T) {
	DebugLog("test", "test")
}

func TestWarnLog(t *testing.T) {
	WarnLog("test", "test")
}

func TestErrorLog(t *testing.T) {
	ErrorLog("test", "test")
}
