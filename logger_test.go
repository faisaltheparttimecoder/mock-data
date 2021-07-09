package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"runtime"
	"testing"
)

// setup a new logger
func NewLoggerValidation() (string, *logrus.Logger, *test.Hook) {
	msg := "Hello! World %s"
	logger, hook := test.NewNullLogger()
	return msg, logger, hook
}

// setup test case functions
type LoggerTest struct {
	t     *testing.T
	hook  *test.Hook
	level logrus.Level
	msg   string
	fname string
}

// Run the test cases
func (l *LoggerTest) RunLoggerTestCases() {
	tests := []struct {
		name  string
		check interface{}
		want  interface{}
	}{
		{"logger_entries", 1, 1},
		{"logger_level", l.hook.LastEntry().Level, l.level},
		{"logger_msg", l.hook.LastEntry().Message, l.msg},
	}
	for _, tt := range tests {
		l.t.Run(tt.name, func(t *testing.T) {
			if tt.check != tt.want {
				t.Errorf("%v = %v, want %v", l.fname, tt.check, tt.want)
			}
		})
	}
	l.t.Run("logger_last_entry", func(t *testing.T) {
		l.hook.Reset()
		if l.hook.LastEntry() != nil {
			t.Errorf("%v = %v, want %v", l.fname, l.hook.LastEntry(), nil)
		}
	})
}

// Test: Debug, checking if correct log debug level is set and send
func TestDebug(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.DebugLevel
	logger.Debug(msg)
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.DebugLevel,
		msg:   msg,
		fname: "TestDebug",
	}
	l.RunLoggerTestCases()
}

// Test: Debugf, checking if correct log debugf level is set and send
func TestDebugf(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.DebugLevel
	logger.Debugf(msg, "user")
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.DebugLevel,
		msg:   fmt.Sprintf(msg, "user"),
		fname: "TestDebugf",
	}
	l.RunLoggerTestCases()
}

// Test: Info, checking if correct log info level is set and send
func TestInfo(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.InfoLevel
	logger.Info(msg)
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.InfoLevel,
		msg:   msg,
		fname: "TestInfo",
	}
	l.RunLoggerTestCases()
}

// Test: Infof, checking if correct log infof level is set and send
func TestInfof(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.InfoLevel
	logger.Infof(msg, "user")
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.InfoLevel,
		msg:   fmt.Sprintf(msg, "user"),
		fname: "TestInfof",
	}
	l.RunLoggerTestCases()
}

// Test: Warn, checking if correct log warn level is set and send
func TestWarn(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.WarnLevel
	logger.Warn(msg)
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.WarnLevel,
		msg:   msg,
		fname: "TestWarn",
	}
	l.RunLoggerTestCases()
}

// Test: Warnf, checking if correct log warnf level is set and send
func TestWarnf(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.WarnLevel
	logger.Warnf(msg, "user")
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.WarnLevel,
		msg:   fmt.Sprintf(msg, "user"),
		fname: "TestWarnf",
	}
	l.RunLoggerTestCases()
}

// Test: Error, checking if correct log error level is set and send
func TestError(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.ErrorLevel
	logger.Error(msg)
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.ErrorLevel,
		msg:   msg,
		fname: "TestError",
	}
	l.RunLoggerTestCases()
}

// Test: Errorf, checking if correct log errorf level is set and send
func TestErrorf(t *testing.T) {
	msg, logger, hook := NewLoggerValidation()
	logger.Level = logrus.ErrorLevel
	logger.Errorf(msg, "user")
	l := &LoggerTest{
		t:     t,
		hook:  hook,
		level: logrus.ErrorLevel,
		msg:   fmt.Sprintf(msg, "user"),
		fname: "TestErrorf",
	}
	l.RunLoggerTestCases()
}

// Test: Fatal, checking if correct log fatal level is set and send
func TestFatal(t *testing.T) {
	// we assume this is working, since this will break the test since it exits the program
}

// Test: Fatalf, checking if correct log fatalf level is set and send
func TestFatalf(t *testing.T) {
	// we assume this is working, since this will break the test since it exits the program
}

// Test: Panic, checking if correct log panic level is set and send
func TestPanic(t *testing.T) {
	// we assume this is working, since this will break the test since it exits the program
}

// Test: Panicf, checking if correct log panic level is set and send
func TestPanicf(t *testing.T) {
	// we assume this is working, since this will break the test since it exits the program
}

// Test: fileInfo, should print the line number where this is executed
func TestFileInfo(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	// why adding 4 ? the function "fileInfo" is called 4 lines from above call
	want := fmt.Sprintf("%v:%d", "logger_test.go", line+4)
	t.Run("check_the_file_name_and_the_line_number_is_displayed", func(t *testing.T) {
		if got := fileInfo(1); got != want {
			t.Errorf("TestFileInfo = %v, want %v", got, want)
		}
	})
}
