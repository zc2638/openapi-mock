package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

/**
 * Created by zc on 2019-11-28.
 */
var (
	Info *logrus.Logger
	Error *logrus.Logger
)

func init() {
	Info = &logrus.Logger{
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		},
		Out:      os.Stderr,
		Level:    logrus.InfoLevel,
		Hooks:    make(logrus.LevelHooks),
		ExitFunc: os.Exit,
	}

	Error = &logrus.Logger{
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		},
		Out:      os.Stderr,
		Level:    logrus.ErrorLevel,
		Hooks:    make(logrus.LevelHooks),
		ExitFunc: os.Exit,
	}

	Info.AddHook(NewHook("info"))
	Error.AddHook(NewHook("error"))
}