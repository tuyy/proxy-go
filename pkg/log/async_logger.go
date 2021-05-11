package log

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const asyncLoggerClassName = "asynclogger"

type AsyncLogger struct {
	level           Level
	formatter       Formatter
	description     string
	baseFilePath    string
	currentFilePath string
	file            *os.File

	itemChan  chan Item
	quitChan  chan interface{}
}

func (l AsyncLogger) Level() Level {
	return l.level
}

func (l AsyncLogger) Description() string {
	return fmt.Sprintf("AsyncLogger:\"%s\"", l.description)
}

func (l *AsyncLogger) Close() {
	for len(l.itemChan) > 0 {
		time.Sleep(time.Millisecond)
	}

	l.quitChan <- true
	l.file.Close()
	l.file = nil
}

func (l *AsyncLogger) Log(item Item) {
	l.itemChan <- item
}

func (l AsyncLogger) writeLog(item Item) {
	filePath := makeCurrentFilePath(l.baseFilePath)
	if filePath != l.currentFilePath {
		l.file.Close()
		newFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		l.currentFilePath = filePath
		l.file = newFile
	}

	fmt.Fprintln(l.file, l.formatter.MakeLog(item))
}

func (l *AsyncLogger) start() {
	go func() {
		for {
			select {
			case item := <-l.itemChan:
				l.writeLog(item)
			case <-l.quitChan:
				return
			}
		}
	}()
}

var emptyFilePathError = errors.New("empty file path in file logger")

func NewAsyncLogger(level Level, formatter Formatter, loggerInfo map[string]string) (Logger, error) {
	logger := &AsyncLogger{
		level:        level,
		formatter:    formatter,
		description:  loggerInfo["description"],
		baseFilePath: loggerInfo["filepath"],
		itemChan:     make(chan Item),
		quitChan:     make(chan interface{}),
	}

	if logger.baseFilePath == "" {
		return nil, emptyFilePathError
	}

	logger.start()

	return logger, nil
}
