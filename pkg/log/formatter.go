package log

import (
	"strconv"
	"strings"
)

type Formatter struct {
	format string
}

const (
	msgFormatStr      = "%m"
	dateTimeFormatStr = "%D"
	pidFormatStr      = "%p"
	tidFormatStr      = "%t"
	fileFormatStr     = "%f"
	lineFormatStr     = "%n"
	levelFormatStr    = "%l"
)

var formatStrList = []string{msgFormatStr, dateTimeFormatStr, tidFormatStr, pidFormatStr, fileFormatStr, lineFormatStr, levelFormatStr}

const logTimeLayout = "2006-01-02 15:04:05.000000"

// TODO 한번만 돌아서 로깅하도록 수정
func (f Formatter) MakeLog(item Item) string {
	result := f.format

	for _, formatStr := range formatStrList {
		switch formatStr {
		case msgFormatStr:
			result = strings.ReplaceAll(result, formatStr, item.msg)
		case dateTimeFormatStr:
			dateTimeStr := item.time.Format(logTimeLayout)
			result = strings.ReplaceAll(result, formatStr, dateTimeStr)
		case pidFormatStr:
			result = strings.ReplaceAll(result, formatStr, strconv.Itoa(item.pid))
		case tidFormatStr:
			result = strings.ReplaceAll(result, formatStr, strconv.Itoa(item.tid))
		case levelFormatStr:
			result = strings.ReplaceAll(result, formatStr, strconv.Itoa(int(item.level)))
		case fileFormatStr:
			result = strings.ReplaceAll(result, formatStr, item.file)
		case lineFormatStr:
			result = strings.ReplaceAll(result, formatStr, strconv.Itoa(item.line))
		}
	}
	return result
}

func NewFormatter(format string) Formatter {
	if format == "" {
		format = msgFormatStr
	}
	return Formatter{
		format: format,
	}
}
