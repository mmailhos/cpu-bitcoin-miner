package main

import (
	"github.com/btcsuite/btcrpcclient"
	"gobtcminer/block"
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

	//Starting a block per thread
	epoch_time := uint32(time.Now().Unix())
	check_chan := make(chan mining.ChannelCheck)
	for i := 0; i < conf.Threads; i++ {
		myblock := block.MakeSemiRandom_BlockHeader(2, epoch_time)
		go mining.Mining_BlockHeader(i, diff, myblock, check_chan)
	}
	for {
		//Mining
	}
}
