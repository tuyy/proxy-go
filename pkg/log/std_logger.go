package log

import "fmt"

const stdLoggerClassName = "stdlogger"

type StdLogger struct {
	level       Level
	formatter   Formatter
	description string
}

func (l StdLogger) Log(item Item) {
	fmt.Println(l.formatter.MakeLog(item))
}

func (l StdLogger) Level() Level {
	return l.level
}

func (l StdLogger) Description() string {
	return fmt.Sprintf("StdLogger:\"%s\"", l.description)
}

func (l StdLogger) Close() {
	// nothing..
}

func NewStdLogger(level Level, formatter Formatter, loggerInfo map[string]string) (*StdLogger, error) {
	return &StdLogger{level: level, formatter: formatter, description: loggerInfo["description"]}, nil
}
