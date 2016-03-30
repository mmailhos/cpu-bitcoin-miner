/*
Author: Mathieu Mailhos
Filename: logger.go
Description: Set the parameters for the logger. It is a home-made logger that also keep tracks of the count of hashes for a easier debugging. The purpose is not to re-write 'log' library but to use it and to include some measurements from the mining pool.
*/

package logger

import (
	"gobtcminer/config"
	"log"
	"sync"
)

//Logger object.
type Logger struct {
	Activated bool       //Can be activatede, or not.
	Level     string     //Levels: 'debug' and 'info'
	File      string     //Filename for storing the output
	HashCount uint32     //Global count of hashes executed so far
	mux       sync.Mutex //Mutex for avoiding concurrency on increasing HashCount
}

//Constructor function
func NewLogger(logger config.JsonLogger) Logger {
	return Logger{
		Activated: logger.Activated,
		Level:     logger.Level,
		File:      logger.File}
}

//Log
func (logger *Logger) Print(level string, output string) {
	if logger.Activated {
		if logger.Level == "info" {
			log.Println(output)
		} else if logger.Level == "debug" {
			log.Println(output)
		}
	}
}
