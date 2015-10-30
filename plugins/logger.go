package plugins

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"strconv"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *log.Logger

func mustInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func mustBool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

func InitializeLogger(filename, namespace string) {
	var f *os.File
	var w io.Writer
	var l *log.Logger
	var err error

	if filename == "" || filename == "stdout" {
		l = log.New(os.Stdout, fmt.Sprint(namespace, " "), log.Ldate|log.Lmicroseconds)
	} else if filename == "null" {
		f, err = os.Open(os.DevNull)
		l = log.New(f, fmt.Sprint(namespace, " "), 0)
	} else if filename == "syslog" {
		w, err = syslog.New(syslog.LOG_NOTICE, namespace)
		l = log.New(w, "", 0)
	} else if strings.HasPrefix(filename, "lumberjack") {
		parts := strings.Split(filename, ",")
		// Expand to 6 parts (defaults are used if parts are missing).
		for len(parts) < 6 {
			parts = append(parts, "")
		}
		l = log.New(&lumberjack.Logger{
			Filename:   parts[1],
			MaxSize:    mustInt(parts[2]),
			MaxAge:     mustInt(parts[3]),
			MaxBackups: mustInt(parts[4]),
			LocalTime:  mustBool(parts[5]),
		}, fmt.Sprint(namespace, " "), log.Ldate|log.Lmicroseconds)
	} else {
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		l = log.New(f, fmt.Sprint(namespace, " "), log.Ldate|log.Lmicroseconds)
	}

	if err != nil {
		panic(errors.New(fmt.Sprint("Couldn't initialize log: ", err)))
	}

	Logger = l
}
