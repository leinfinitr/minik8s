package log

import (
	"testing"
)

func TestInfoLog(t *testing.T) {
	InfoLog("test")
}

func TestDebugLog(t *testing.T) {
	DebugLog("test")
}

func TestWarnLog(t *testing.T) {
	WarnLog("test")
}

func TestErrorLog(t *testing.T) {
	ErrorLog("test")
}
