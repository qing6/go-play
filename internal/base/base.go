package base

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

// 判断当先是否为"go test", 这个方法貌似并不是百分百靠谱
// https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
func InTesting() bool {
	return flag.Lookup("test.v") != nil
}

func Fail(msg string, fields Fields) error {
	return errors.New(gainErrorMsg(msg, fields))
}

func FailBy(err error, msg string, fields Fields) error {
	return errors.Wrap(err, gainErrorMsg(msg, fields))
}

func gainErrorMsg(msg string, fields Fields) string {
	buf := BufPool.Get().(*bytes.Buffer)
	defer BufPool.Put(buf)
	buf.Reset()
	buf.WriteString(msg)
	fields.exportTo(buf)
	return buf.String()
}

type Fields map[string]interface{}

func (fields Fields) exportTo(buf *bytes.Buffer) {
	for k, v := range fields {
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(fmt.Sprint(v))
	}
}

var BufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	if InTesting() {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetOutput(os.Stdout)
	} else {
		execDir, execFilename, _ := SplitFilepath(os.Args[0])
		logDir := filepath.Join(execDir, fmt.Sprintf("%s-log", execFilename))
		os.Mkdir(logDir, 0666)
		logFile := filepath.Join(logDir, fmt.Sprintf("%sstart.log", time.Now().Format(time.RFC3339)))
		logger.SetOutput(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}
	return logger
}

func SplitFilepath(file string) (dir, filename, fileExt string) {
	dir, filename = filepath.Split(file)
	fileExt = filepath.Ext(filename)
	filename = strings.TrimRight(filename, fileExt)
	return
}
