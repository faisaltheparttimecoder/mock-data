package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// SetLogLevel set the logger level
func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

// SetLogFormatter set the log entry format
func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Debug(args...)
	}
}

// Debugf logs with format message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Debugf(format, args...)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Info(args...)
	}
}

// Infof logs with format message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Infof(format, args...)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Warn(args...)
	}
}

// Warnf logs with format message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Warnf(format, args...)
	}
}

// Error logs a message at level Error on the standard logger. [unused]
//func Error(args ...interface{}) {
//	if logger.Level >= logrus.ErrorLevel {
//		entry := logger.WithFields(logrus.Fields{})
//		if cmdOptions.Debug {
//			entry.Data["file"] = fileInfo(2)
//		}
//		entry.Error(args...)
//	}
//}

// Errorf logs with format message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Errorf(format, args...)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Fatal(args...)
	}
}

// Fatalf logs with format message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		if cmdOptions.Debug {
			entry.Data["file"] = fileInfo(2)
		}
		entry.Fatalf(format, args...)
	}
}

// Panic logs a message at level Panic on the standard logger.[unused]
//func Panic(args ...interface{}) {
//	if logger.Level >= logrus.PanicLevel {
//		entry := logger.WithFields(logrus.Fields{})
//		if cmdOptions.Debug {
//			entry.Data["file"] = fileInfo(2)
//		}
//		entry.Panic(args...)
//	}
//}

// Panic logs with format message at level Panic on the standard logger.[unused]
//func Panicf(format string, args ...interface{}) {
//	if logger.Level >= logrus.PanicLevel {
//		entry := logger.WithFields(logrus.Fields{})
//		if cmdOptions.Debug {
//			entry.Data["file"] = fileInfo(2)
//		}
//		entry.Panicf(format, args...)
//	}
//}

// Display the file and the line number where the command
// was executed
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// Setup or Initialize the logger
func initLogger(verbose bool) {
	// Set the formatter option for logrus
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.FullTimestamp = true

	// Set the formatter.
	SetLogFormatter(formatter)

	// if log level
	if verbose {
		SetLogLevel(logrus.DebugLevel)
	} else {
		SetLogLevel(logrus.InfoLevel)
	}
}
