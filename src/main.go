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
	dispatcher := mining.NewDispatcher(monitor)
	filling := true
	//Run new chunks in the jobqueue
	for filling == true {
		for _, chunk := range mining.NewChunkList(2, uint32(time.Now().Unix()), diff) {
			if len(dispatcher.ChunkQueue) < cap(dispatcher.ChunkQueue) {
				dispatcher.ChunkQueue <- chunk
			} else {
				filling = false
				break
			}
		}
	}
	dispatcher.Run()
}
