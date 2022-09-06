package utils

import (
	"log"
	"os"
)

var InfoLog = log.New(os.Stdout, "INFO:\t", log.LstdFlags)
var ErrorLog = log.New(os.Stdout, "ERROR:\t", log.LstdFlags|log.Llongfile)
