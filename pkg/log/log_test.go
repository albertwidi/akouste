package log

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLineNumber(t *testing.T) {
	SetConfig(&Config{"debug", "test.log", "test.log", true})

	Debug("test")
	Debugf("test")
	Debugw("test", nil)

	Print("test")
	Println("test")
	Printf("test")
	Info("test")
	Infof("test")
	Infow("test", nil)

	Warn("test")
	Warnf("test")
	Warnw("test", nil)

	Error("test")
	Errorf("test")
	Errorw("test", nil)

	// Fatal("test")
	// Fatalf("test")
	// Fatalw("test", nil)

	text, err := ioutil.ReadFile("test.log")
	if err != nil {
		t.Error(err)
	}

	caller := strings.Count(string(text), "log_test.go")
	expected := 15
	if caller != expected {
		t.Errorf("Wrapper is not uniform, expected %d, got %d", expected, caller)
	}
	os.Remove("test.log")
}
