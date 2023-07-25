package logger

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type LogLevel string

const (
	DEBUG   LogLevel = "DEBUG"
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
)

type Logger struct {
	level   LogLevel
	writeTo io.Writer
}

func NewLogger(level LogLevel, writeTo io.Writer) *Logger {
	return &Logger{level, writeTo}
}

func (l Logger) msg(level LogLevel, template string, a ...any) {
	var buildedString strings.Builder
	buildedString.WriteString(fmt.Sprintf("%s [%s] ", level, time.Now().UTC().Format("2006-01-02 15:04:05")))
	buildedString.WriteString(fmt.Sprintf(template, a...))
	if !strings.HasSuffix(template, "\n") {
		buildedString.WriteString("\n")
	}
	l.writeTo.Write([]byte(buildedString.String()))
}

func (l Logger) Debug(template string, a ...any) {
	if l.level == DEBUG {
		l.msg(DEBUG, template, a...)
	}
}

func (l Logger) Info(template string, a ...any) {
	if l.level == INFO || l.level == DEBUG {
		l.msg(INFO, template, a...)
	}
}

func (l Logger) Warning(template string, a ...any) {
	if l.level == WARNING || l.level == INFO || l.level == DEBUG {
		l.msg(WARNING, template, a...)
	}
}

func (l Logger) Error(template string, a ...any) {
	if l.level == ERROR || l.level == WARNING || l.level == INFO || l.level == DEBUG {
		l.msg(ERROR, template, a...)
	}
}

func (l Logger) Log(template string, a ...any) {
	l.msg(l.level, template, a...)
}
