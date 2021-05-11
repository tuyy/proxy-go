package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Level int

const (
	debug Level = iota + 1
	info
	err
	fatal
	access
	result
)

type Logger interface {
	Close()
	Description() string
	Level() Level
	Log(item Item)
}

type Item struct {
	time           time.Time
	level          Level
	file, msg      string
	line, pid, tid int
}

type Container struct {
	loggers map[Level]Logger
}

var logContainer *Container

func Init(loggerInfoList []map[string]string) {
	logContainer = &Container{loggers: make(map[Level]Logger)}

	for _, loggerInfo := range loggerInfoList {
		logger, err := getLogger(loggerInfo)
		if err != nil {
			panic(fmt.Sprintf("invalid logger info. err:%s info:%+v", err, loggerInfo))
		}
		logContainer.registerLogger(logger.Level(), logger)
	}
}

func getLogger(loggerInfo map[string]string) (Logger, error) {
	levelNo, err := strconv.ParseInt(loggerInfo["level"], 10, 32)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid logger level. err:%s", err))
	}
	level := Level(levelNo)
	formatter := NewFormatter(loggerInfo["format"])

	switch strings.ToLower(loggerInfo["class"]) {
	case stdLoggerClassName:
		return NewStdLogger(level, formatter, loggerInfo)
	case fileLoggerClassName:
		return NewFileLogger(level, formatter, loggerInfo)
	case asyncLoggerClassName:
		return NewAsyncLogger(level, formatter, loggerInfo)
	default:
		return nil, errors.New(fmt.Sprintf("invalid logger class name. class:%s", loggerInfo["class"]))
	}
}

func (c *Container) registerLogger(level Level, logger Logger) {
	c.loggers[level] = logger
}

func (c Container) log(level Level, msg string) {
	_, file, line, _ := runtime.Caller(2)
	item := Item{
		time: time.Now(),
		level: level,
		msg: msg,
		file: file,
		line: line,
		pid: os.Getpid(),
		tid: os.Getpid(),
	}

	if logger, ok := c.loggers[level]; ok {
		logger.Log(item)
	} else {
		fmt.Println(msg)
	}
}

func (c Container) close() {
	for _, logger := range c.loggers {
		logger.Close()
	}
}

func Custom(level Level, msg string) {
	if logContainer == nil {
		fmt.Println(msg)
	} else {
		logContainer.log(level, msg)
	}
}

func Debug(format string, a ...interface{}) {
	Custom(debug, fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) {
	Custom(info, fmt.Sprintf(format, a...))
}

func Error(format string, a ...interface{}) {
	Custom(err, fmt.Sprintf(format, a...))
}

func Fatal(format string, a ...interface{}) {
	Custom(fatal, fmt.Sprintf(format, a...))
}

func Access(format string, a ...interface{}) {
	Custom(access, fmt.Sprintf(format, a...))
}

func Result(format string, a ...interface{}) {
	Custom(result, fmt.Sprintf(format, a...))
}
