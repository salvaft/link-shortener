package utils

import (
	"log"
	"os"
)

func getLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

var Logger = getLogger()
