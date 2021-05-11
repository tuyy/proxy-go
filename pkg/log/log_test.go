package log

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestInitHappy(t *testing.T) {
	loggerInfoList := []map[string]string{
		{"level": "1", "description": "debug", "class": "StdLogger"},
		{"level": "2", "description": "info", "filepath": "myLog1.log", "class": "FileLogger"},
		{"level": "99", "description": "result", "filepath": "myLog2.log", "class": "AsyncLogger"},
	}

	Init(loggerInfoList)

	expected := 3
	result := len(logContainer.loggers)
	assert.Equal(t, expected, result)

	logContainer.close()
}

func TestGetLoggerHappy(t *testing.T) {
	cases := []struct {
		input map[string]string
		want  string
	}{
		{
			input: map[string]string{"level": "1", "description": "debug", "class": "StdLogger"},
			want:  "StdLogger:\"debug\"",
		},
		{
			input: map[string]string{"level": "2", "description": "info", "filepath": "myLog1.log", "class": "FileLogger"},
			want:  "FileLogger:\"info\"",
		},
		{
			input: map[string]string{"level": "4", "description": "fatal", "filepath": "myLog2.log", "class": "AsyncLogger"},
			want:  "AsyncLogger:\"fatal\"",
		},
	}

	for _, c := range cases {
		result, _ := getLogger(c.input)
		assert.Equal(t, c.want, result.Description())
	}
}

func TestGetLoggerWhenLevelIsInvalid(t *testing.T) {
	cases := []struct {
		input map[string]string
	}{
		{input: map[string]string{"level": "", "description": "debug", "class": "StdLogger"}},
		{input: map[string]string{"level": "hello", "description": "debug", "class": "StdLogger"}},
		{input: map[string]string{"level": "99999999999", "description": "debug", "class": "StdLogger"}},
	}

	for _, c := range cases {
		_, err := getLogger(c.input)
		assert.NotNil(t, err)
	}
}

func TestGetLoggerWhenBaseFilePathIsEmpty(t *testing.T) {
	cases := []struct {
		input map[string]string
	}{
		{input: map[string]string{"level": "2", "description": "info", "filepath": "", "class": "AsyncLogger"}},
		{input: map[string]string{"level": "3", "description": "err", "class": "AsyncLogger"}},
		{input: map[string]string{"level": "2", "description": "info", "filepath": "", "class": "FileLogger"}},
		{input: map[string]string{"level": "3", "description": "err", "class": "FileLogger"}},
	}

	for _, c := range cases {
		_, err := getLogger(c.input)
		assert.Equal(t, emptyFilePathError, err)
	}
}

func TestWriteFileHappy(t *testing.T) {
	cases := []struct {
		class      string
		loggerFunc func(Level, Formatter, map[string]string) (Logger, error)
		want       string
	}{
		{class: "FileLogger", loggerFunc: NewFileLogger, want: "1\n2\n3\n"},
		{class: "AsyncLogger", loggerFunc: NewAsyncLogger, want: "1\n2\n3\n"},
	}

	baseFilePath := "myLog2.log"
	for _, c := range cases {
		loggerInfo := map[string]string{"level": "1", "description": "debug", "filepath": baseFilePath, "class": c.class}
		logger, _ := c.loggerFunc(debug, NewFormatter(""), loggerInfo)

		filePath := makeCurrentFilePath(baseFilePath)
		os.Remove(filePath)

		logger.Log(Item{msg: "1"})
		logger.Log(Item{msg: "2"})
		logger.Log(Item{msg: "3"})

		b, _ := ioutil.ReadFile(filePath)
		expected := "1\n2\n3\n"
		result := string(b)
		assert.Equal(t, expected, result)

		logger.Close()
		os.Remove(filePath)
	}
}

func TestFormatterHappy(t *testing.T) {
	dateTime := time.Date(2021, 2, 1, 0, 0, 0, 0, time.Local)
	item := Item{
		time:  dateTime,
		level: debug,
		msg:   "hello world",
		file:  "test.go",
		line:  10,
		pid:   101,
		tid:   1002,
	}

	cases := []struct {
		input string
		want  string
	}{
		{
			input: "",
			want:  "hello world",
		},
		{
			input: "hello",
			want:  "hello",
		},
		{
			input: "%m",
			want:  "hello world",
		},
		{
			input: "[%D] %p:%t %l %f:%n %m",
			want:  "[2021-02-01 00:00:00.000000] 101:1002 1 test.go:10 hello world",
		},
	}

	for _, c := range cases {
		formatter := NewFormatter(c.input)
		result := formatter.MakeLog(item)
		assert.Equal(t, c.want, result)
	}
}
