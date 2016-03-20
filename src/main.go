package main

import (
	"encoding/json"
	"github.com/btcsuite/btcrpcclient"
	"io/ioutil"
	"log"
	"os/exec"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Account  string `json:"account"`
}

//Missing: capabilities, transactions, mutable
type ResultTemplate struct {
	PreviousBlockHash string `json:"previousblockhash"`
	Target            string `json:"target"`
	NonceRange        string `json:"noncerange"`
	MinTime           uint   `json:"mintime"`
	SigOpLimit        uint   `json:"sigoplimit"`
	CurTime           uint   `json:"curtime"`
	Height            uint   `json:"height"`
	Version           uint   `json:"version"`
	Bits              string `json:"bits"`
	CoinBaseValue     uint   `json:"coinbasevalue"`
	SizeLimit         uint   `json:"sizelimit"`
	LongPollId        string `json:"longpollid"`
}

type BlockTemplate struct {
	Error  string         `json:"error"`
	Result ResultTemplate `json:"result"`
}

func VerifyAccount(client *btcrpcclient.Client, name string) bool {
	adr, err := client.GetAccountAddress(name)
	if err != nil {
		log.Printf("Error getting account address %s", name)
		return false
	} else {
		wal, err := client.ValidateAddress(adr)
		if err != nil {
			log.Printf("Error validating account address")
			return false
		} else if !wal.IsValid {
			log.Printf("Invalid account address")
			return false
		}
		log.Printf("Account: %s, Address: %s, PubKey: %s\n", name, adr, wal.PubKey)
	}
	return true
}
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

func readconf() (conf Config) {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error:", err)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("Error:", err)
	}
	return
}

// VERY Temporary work-around for GetBlockTemplate ;)
func GetBlockTemplate(user, password, host string) ResultTemplate {
	command := "curl -u " + user + ":" + password + ` --data-binary '{"jsonrpc": "1.0", "id":"curltest", "method": "getblocktemplate", "params": [] }'   -H 'content-type: text/plain;' http://` + host + "/"
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}
	var btp BlockTemplate
	err = json.Unmarshal(out, &btp)
	if err != nil {
		log.Fatalf("Error:", err)
	}
	return btp.Result
}

func main() {
	// Read and parse the configuration file
	conf := readconf()
	// Create new client instance
	client, err := btcrpcclient.New(&btcrpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Password,
	}, nil)
	if err != nil {
		log.Fatalf("error creating new btc client: %v", err)
	}
	// Verifying Account
	if !VerifyAccount(client, conf.Account) {
		ListAccounts(client)
	}
	// Get and Parse BlockTemplate to begin minning
	GetBlockTemplate(conf.User, conf.Password, conf.Host)
}
