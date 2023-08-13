package logger

import "github.com/sirupsen/logrus"

type DebugWriter struct {
	Logger *logrus.Logger
}

func (w *DebugWriter) Write(p []byte) (int, error) {
	w.Logger.Debug(string(p))
	return len(p), nil
}
