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
	"time"
)

//Logger object.
type Logger struct {
	Activated  bool       //Can be activatede, or not.
	Level      string     //Levels: 'debug' and 'info'
	File       string     //Filename for storing the output
	HashCount  uint32     //Global count of hashes executed so far
	BlockCount uint32     //Global count of hashes executed so far
	mux        sync.Mutex //Mutex for avoiding concurrency on increasing HashCount
	BeginTime  time.Time  //Used for calculating the compute time for benchmarking
}

//NewLogger Constructor function
func NewLogger(logger config.JSONLogger) Logger {
	return Logger{
		Activated: logger.Activated,
		Level:     logger.Level,
		File:      logger.File,
	}
}

//Print simply logs
func (logger *Logger) Print(level string, output string) {
	if logger.Activated {
		if logger.Level == level {
			log.Println(output)
		} else if logger.Level == "debug" && level == "info" {
			log.Println(output)
		}
	}
}

//IncrementHashCount counts the number of hash executed regularly. Use of a mutex to avoid race condition.
func (logger *Logger) IncrementHashCount(count uint32) {
	logger.mux.Lock()
	defer logger.mux.Unlock()
	logger.HashCount += count
}

//IncrementBlockCount increments the number of succesfuly mined block. Use of a mutex to avoid race condition.
func (logger *Logger) IncrementBlockCount() {
	logger.mux.Lock()
	defer logger.mux.Unlock()
	logger.BlockCount++
}
