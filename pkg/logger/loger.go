package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type fileHook struct {
	LevelsArr []logrus.Level
	Files     map[logrus.Level]*os.File
}

func (hook *fileHook) Fire(entry *logrus.Entry) error {
	for _, level := range hook.LevelsArr {
		if entry.Level <= level {
			entry.Logger.Out = hook.Files[level]
			break
		}
	}
	return nil
}

func (hook *fileHook) Levels() []logrus.Level {
	return hook.LevelsArr
}

func InitLogger() {
	logger = logrus.New()

	debugFile, err := os.OpenFile(
		"pkg/logger/logs/debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		logrus.Fatalf("Failed to open debug log file: %v", err)
	}

	fatalFile, err := os.OpenFile(
		"pkg/logger/logs/fatal.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		logrus.Fatalf("Failed to open fatal log file: %v", err)
	}

	logger.AddHook(&fileHook{
		LevelsArr: []logrus.Level{
			logrus.FatalLevel,
			logrus.DebugLevel,
		},
		Files: map[logrus.Level]*os.File{
			logrus.FatalLevel: fatalFile,
			logrus.DebugLevel: debugFile,
		},
	})
	logger.SetLevel(logrus.DebugLevel)
}

func CloseLoggerFile() {
	fileHook, ok := logger.Hooks[logrus.FatalLevel][0].(*fileHook)
	if ok {
		for _, file := range fileHook.Files {
			err := file.Close()
			if err != nil {
				logger.Errorf("Failed to close log file: %v", err)
			}
		}
	}
}

func GetLogger() *logrus.Logger {
	return logger
}
