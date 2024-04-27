package log

import (
	"testing"
)

func TestPrintln(t *testing.T) {
	SetLog(true)
	Println("open test")
	SetLog(false)
	Println("close test")
}
