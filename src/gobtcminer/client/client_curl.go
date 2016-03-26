/*
Author: Mathieu Mailhos
Filename: client_curl.go
Description: Temporary functions used to run request on bitcoin client. Depending on the open-source btcrpcclient project, this file will be overwritten by client_lib.go or re-made from scratch using proper HTTP client.
*/
package client

import (
	"encoding/json"
	"errors"
	"os/exec"
)

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

type Difficulty struct {
	Error      string  `json:"error"`
	Difficulty float64 `json:"result"`
	Id         string  `json:"id"`
}

// VERY Temporary work-around for GetBlockTemplate() from BP023 ;)
//GetResultTemplate(user, password, host)
//Get and parse data received from Bitcoin client on a getblocktemplate request
func GetResultTemplate(user, password, host string) (rtp ResultTemplate, err error) {
	var btp BlockTemplate
	command := "curl -u " + user + ":" + password + ` --data-binary '{"jsonrpc": "1.1", "id":"0", "method": "getblocktemplate", "params": [{"capabilities": ["coinbasetxn", "workid", "coinbase/append"]}] }'   -H 'content-type: application/json;' http://` + host + "/"
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
	}
	return btp.Result, errors.New(btp.Error)
}

func GetDifficulty(user, password, host string) (difficulty float64, err error) {
	command := "curl -u " + user + ":" + password + ` --data-binary '{"jsonrpc": "1.1", "id":"0", "method": "getdifficulty"}'   -H 'content-type: application/json;' http://` + host + "/"
	var dif Difficulty
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return
	}
	err = json.Unmarshal(out, &dif)
	if err != nil {
		return
	}
	if dif.Error == "" {
		return dif.Difficulty, nil
	}
	return dif.Difficulty, errors.New(dif.Error)
}
