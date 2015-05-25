package plugins

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
)

var Logger *log.Logger

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
	} else {
		f, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		l = log.New(f, fmt.Sprint(namespace, " "), log.Ldate|log.Lmicroseconds)
	}

	if err != nil {
		panic(errors.New(fmt.Sprint("Couldn't initialize log: ", err)))
	}

	Logger = l
}
