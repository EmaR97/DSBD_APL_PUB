package utility

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	once sync.Once

	warningLog *log.Logger
	infoLog    *log.Logger
	errorLog   *log.Logger
)

func initializeLogger() {
	file, err := os.OpenFile("myLOG.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Create a multi-writer to log to both file and console
	multiWriter := io.MultiWriter(file, os.Stdout)

	infoLog = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(multiWriter, "\033[33mWARNING\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(multiWriter, "\033[31mERROR\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// InfoLog returns the InfoLog instance
func InfoLog() *log.Logger {
	once.Do(initializeLogger)
	return infoLog
}

// WarningLog returns the WarningLog instance
func WarningLog() *log.Logger {
	once.Do(initializeLogger)
	return warningLog
}

// ErrorLog returns the ErrorLog instance
func ErrorLog() *log.Logger {
	once.Do(initializeLogger)
	return errorLog
}
