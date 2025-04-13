package logger

import (
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"io"
)

func LogRotator(path string) io.Writer {
	logger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1,    //days
		Compress:   true, // disabled by default
	}
	return logger
}
