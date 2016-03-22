/*
Author: Mathieu Mailhos
Filename: client_lib.go
Description: Use of btcrpcclient to set the environment and get the block header for starting mining. On-hold for now since important functions are currently being developed on the project due to recent BP0023 changes. See: https://github.com/btcsuite/btcrpcclient
*/
package client

import (
	"github.com/btcsuite/btcrpcclient"
	"log"
)

//VerifyAccount(client, name)
//Validates an account name on a rpcclient and check if it has a valid wallet
func VerifyAccount(client *btcrpcclient.Client, name string) (bool, error) {
	adr, err := client.GetAccountAddress(name)
	if err != nil {
		return false, err
	} else {
		wal, err := client.ValidateAddress(adr)
		if err != nil {
			return false, err
		} else if !wal.IsValid {
			return false, err
		}
	}
	return true, nil
}

//ListAccounts(client)
//Informative function listing accounts available
func ListAccounts(client *btcrpcclient.Client) {
	accounts, err := client.ListAccounts()
	if err != nil {
		log.Fatalf("Error listing accounts: %v", err)
	}
	for label, amount := range accounts {
		log.Println("Account %s with %s", label, amount)
	}
	log.Fatalf("Indicates the right account in config.json then try again.")
}
