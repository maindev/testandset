package util

import (
	"log"
	"os"
)

// Verbose global verbose logging flag
var Verbose bool

// SetUpLogging set up logging
func SetUpLogging() {
	log.SetOutput(os.Stderr)
}

// WriteVerboseMessage write log message when verbose flag enabled
func WriteVerboseMessage(message string) {
	if Verbose {
		log.Println(message)
	}
}

// ExitWithMessage log message and exit with exit code 1
func ExitWithMessage(message string) {
	log.Println(message)
	os.Exit(1)
}
