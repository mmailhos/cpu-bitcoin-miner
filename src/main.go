package main

import (
	"github.com/btcsuite/btcrpcclient"
	"gobtcminer/client"
	"gobtcminer/config"
	"log"
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
	//Loading and parsing values from Bitcoin API call
	_, err = client.GetResultTemplate(conf.User, conf.Password, conf.Host)
	if err != nil {
		log.Fatalf("Error getting mining data: %v", err)
	}
}
