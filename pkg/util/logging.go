package util

import (
	"github.com/sirupsen/logrus"
)

func LogHandler(handler string, msg string, err error) {
	logrus.WithFields(logrus.Fields{
		"handler": handler,
		"msg":     msg,
	}).Error(err)
}
