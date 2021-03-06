package rest

import (
	_ "io"
	_ "io/ioutil"

	"log"
	"runtime"
	"time"

	_ "gopkg.in/natefinch/lumberjack.v2"
)

type Location struct {
	file string
	line int
}

var (
	TScript *log.Logger
	TLog    *log.Logger
	TTime   *log.Logger
	TResult *log.Logger
)

func initLog(bScript, bTime bool) {
	TLog = api.logger
	TTime = api.logger
	TScript = api.logger
	TResult = api.logger

	/*
		// General loger
		gLog := &lumberjack.Logger{Filename: "EkaraApi.log", MaxSize: 250, MaxBackups: 5, MaxAge: 5}
		TLog = log.New(gLog, "Log: ", log.Ldate|log.Ltime|log.Lshortfile)
		TResult = log.New(gLog, "Result: ", log.Ldate|log.Ltime|log.Lshortfile)

		var writer io.Writer

		// Time execution loger
		TLog.Printf("time loggin required %v", bTime)
		if bTime {
			writer = &lumberjack.Logger{Filename: "EkaraApiTime.log", MaxSize: 250, MaxBackups: 5, MaxAge: 5}
		} else {
			writer = ioutil.Discard // If the log of time is not wanted --> redirect
		}
		TTime = log.New(writer, "Time: ", log.Ldate|log.Ltime|log.Lshortfile)

		// Script logger
		TLog.Printf("script loggin required %v", bScript)
		if bScript {
			writer = &lumberjack.Logger{Filename: "EkaraApiScript.log", MaxSize: 250, MaxBackups: 5, MaxAge: 5}
		} else {
			writer = ioutil.Discard // If the log of script is not wanted --> redirect
		}
		TScript = log.New(writer, "Script: ", log.Ldate|log.Ltime|log.Lshortfile)
		TLog.Printf(LOGGER_INITIALIZED)
	*/
}

func traceTime(location Location) func() {
	t := time.Now()
	return func() {
		TTime.Printf(TIME_REPORT, location.file, location.line, (time.Since(t)))
	}
}

func here() Location {
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	return Location{file, line}
}
