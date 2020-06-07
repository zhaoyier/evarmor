package log

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

// git.ezbuy.me/ezbuy/vendor/github.com/sirupsen/logrus
// git.ezbuy.me/ezbuy/vendor/github.com/Sirupsen/logrus

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func NewFields(fields map[string]interface{}) *logrus.Entry {
	return log.WithFields(fields)
}

func Infof(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{}).Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.WithFields(logrus.Fields{}).Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) error {
	log.WithFields(logrus.Fields{}).Errorf(format, args...)
	return fmt.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) error {
	log.WithFields(logrus.Fields{}).Fatalf(format, args...)
	return fmt.Errorf(format, args...)
}
