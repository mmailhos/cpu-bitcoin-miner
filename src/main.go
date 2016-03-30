package main

import (
	"github.com/btcsuite/btcrpcclient"
	"gobtcminer/client"
	"gobtcminer/config"
	"gobtcminer/mining"
	"log"
	"time"
)

func main() {
	// Read and parse the configuration file
	conf := config.Readconf("config.json")
	// Create new client instance
	rpcclient, err := btcrpcclient.New(&btcrpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Password,
	}, nil)
	if err != nil {
		log.Fatalf("Error creating new btc client: %v", err)
	}

	// Verifying Account
	if val, err := client.VerifyAccount(rpcclient, conf.Account); !val {
		log.Printf("Error: %v ", err)
		client.ListAccounts(rpcclient)
	}

	diff, err := client.GetDifficulty(conf.User, conf.Password, conf.Host)
	if err != nil {
		log.Fatal("Error getting difficulty: %v", err)
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
