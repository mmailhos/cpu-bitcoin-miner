package main

import (
	"encoding/json"
	"errors"
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

// Missing: depends[]
type TransactionTemplate struct {
	Hash   string `json:"hash"`
	Fee    uint   `json:"fee"`
	Data   string `json:"data"`
	SigOps uint   `json:"sigops"`
}

//Missing: capabilities, mutable
type ResultTemplate struct {
	PreviousBlockHash string                `json:"previousblockhash"`
	Target            string                `json:"target"`
	NonceRange        string                `json:"noncerange"`
	Bits              string                `json:"bits"`
	LongPollId        string                `json:"longpollid"`
	MinTime           uint                  `json:"mintime"`
	SigOpLimit        uint                  `json:"sigoplimit"`
	CurTime           uint                  `json:"curtime"`
	Height            uint                  `json:"height"`
	Version           uint                  `json:"version"`
	CoinBaseValue     uint                  `json:"coinbasevalue"`
	SizeLimit         uint                  `json:"sizelimit"`
	Transactions      []TransactionTemplate `json:"transactions"`
}

type BlockTemplate struct {
	Error  string         `json:"error"`
	Result ResultTemplate `json:"result"`
}

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

// VERY Temporary work-around for GetBlockTemplate() from BP023 ;)
func GetResultTemplate(user, password, host string) (rtp ResultTemplate, err error) {
	var btp BlockTemplate
	command := "curl -u " + user + ":" + password + ` --data-binary '{"jsonrpc": "1.1", "id":"0", "method": "getblocktemplate", "params": [{"capabilities": ["coinbasetxn", "workid", "coinbase/append"]}] }'   -H 'content-type: application/json;' http://` + host + "/"
	log.Printf(command)
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return
	}
	err = json.Unmarshal(out, &btp)
	if err != nil {
		return
	}
	if btp.Error == "" {
		return btp.Result, nil
	} else {
		return rtp, errors.New(btp.Error)
	}
	return
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
		log.Fatalf("Error creating new btc client: %v", err)
	}
	// Verifying Account
	if val, err := VerifyAccount(client, conf.Account); !val {
		log.Printf("Error: %v ", err)
		ListAccounts(client)
	}
	//Loading and parsing values from Bitcoin API call
	_, err = GetResultTemplate(conf.User, conf.Password, conf.Host)
	if err != nil {
		log.Fatalf("Error getting mining data: %v", err)
	}
}
