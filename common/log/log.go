package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
)

// git.ezbuy.me/ezbuy/vendor/github.com/sirupsen/logrus
// git.ezbuy.me/ezbuy/vendor/github.com/Sirupsen/logrus

var logger *logrus.Entry
var initLogrusOnce sync.Once

func init() {
	initLogrusOnce.Do(func() {
		log := logrus.New()
		// 为当前logrus实例设置消息的输出,同样地,
		// 可以设置logrus实例的输出到任意io.writer
		log.Out = os.Stdout
		disableColor := os.Getenv("TERM") == "dumb"
		if !disableColor {
			disableColor = os.Getenv("NOCOLOR") != ""
		}
		// 为当前logrus实例设置消息输出格式为json格式.
		// 同样地,也可以单独为某个logrus实例设置日志级别和hook,这里不详细叙述.
		// log.Formatter = &logrus.JSONFormatter{}
		// 设置日志等级
		log.SetLevel(logrus.WarnLevel)
		log.Formatter = &logrus.TextFormatter{
			DisableColors: disableColor,
			SortingFunc:   sortingLogger,
		}
		logger = logrus.NewEntry(log)
	})
}

func sortingLogger(keys []string) {
	keysLastIdx := len(keys) - 1
	for i, k := range keys {
		if k == "msg" && i != keysLastIdx {
			keys[i], keys[keysLastIdx] = keys[keysLastIdx], keys[i]
			break
		}
	}
}

func NewFields(fields map[string]interface{}) *logrus.Entry {
	return logger.WithFields(fields)
}

func Infof(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) error {
	logger.WithFields(logrus.Fields{}).Errorf(format, args...)
	return fmt.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) error {
	logger.WithFields(logrus.Fields{}).Fatalf(format, args...)
	return fmt.Errorf(format, args...)
}
