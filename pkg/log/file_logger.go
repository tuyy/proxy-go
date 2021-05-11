package log

import (
	"fmt"
	"os"
	"time"
)

const fileLoggerClassName = "filelogger"

type FileLogger struct {
	level           Level
	formatter       Formatter
	description     string
	baseFilePath    string
	currentFilePath string
	file            *os.File
}

func (f FileLogger) Level() Level {
	return f.level
}

func (f FileLogger) Description() string {
	return fmt.Sprintf("FileLogger:\"%s\"", f.description)
}

func (f *FileLogger) Close() {
	f.file.Close()
	f.file = nil
}

func (f *FileLogger) Log(item Item) {
	filePath := makeCurrentFilePath(f.baseFilePath)
	if filePath != f.currentFilePath {
		f.file.Close()
		newFile, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		f.currentFilePath = filePath
		f.file = newFile
	}

	fmt.Fprintln(f.file, f.formatter.MakeLog(item))
}

const dateTimeLayout = "20060102"

func makeCurrentFilePath(baseFilePath string) string {
	return fmt.Sprintf("%s.%s", baseFilePath, time.Now().Format(dateTimeLayout))
}

// TODO marshal 을 할까?
func NewFileLogger(level Level, formatter Formatter, loggerInfo map[string]string) (Logger, error) {
	logger := &FileLogger{
		level:        level,
		formatter:    formatter,
		description:  loggerInfo["description"],
		baseFilePath: loggerInfo["filepath"],
	}

	if logger.baseFilePath == "" {
		return nil, emptyFilePathError
	}

	return logger, nil
}
