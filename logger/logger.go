package logger

import (
	"flag"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	// set location of log file
	var logpath = "/default.log"
	flag.Parse()
	var file, err1 = os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
