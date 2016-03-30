package main

import (
	"gobtcminer/client"
	"gobtcminer/config"
	"gobtcminer/logger"
	"gobtcminer/mining"
	"time"
)

var monitor logger.Logger

func main() {
	// Read and parse the configuration file
	conf := config.Readconf("config.json")
	monitor = logger.NewLogger(conf.Log)
	diff, err := client.GetDifficulty(conf.User, conf.Password, conf.Host)
	if err != nil {
		monitor.Print("info", "Error getting difficulty: "+err.Error())
	}

	epoch_time := uint32(time.Now().Unix())
	dispatcher := mining.NewDispatcher(conf.Threads)

	//Run new chunks in the jobqueue
	for i := 0; i < 30; i++ {
		if len(dispatcher.ChunkQueue) < cap(dispatcher.ChunkQueue) {
			dispatcher.ChunkQueue <- mining.NewChunk(2, epoch_time, diff)
		}
	}
	dispatcher.Run()
}
