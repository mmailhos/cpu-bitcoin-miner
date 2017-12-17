/*Package client for interacting with bitcoin client
Author: Mathieu Mailhos
Filename: client_curl.go
*/
package client

import (
	"encoding/json"
	"errors"
	"os/exec"
)

//TransactionTemplate Missing: depends[]
type TransactionTemplate struct {
	Hash   string `json:"hash"`
	Fee    uint   `json:"fee"`
	Data   string `json:"data"`
	SigOps uint   `json:"sigops"`
}

//ResultTemplate Missing: capabilities, mutable
type ResultTemplate struct {
	PreviousBlockHash string                `json:"previousblockhash"`
	Target            string                `json:"target"`
	NonceRange        string                `json:"noncerange"`
	Bits              string                `json:"bits"`
	LongPollID        string                `json:"longpollid"`
	MinTime           uint                  `json:"mintime"`
	SigOpLimit        uint                  `json:"sigoplimit"`
	CurTime           uint                  `json:"curtime"`
	Height            uint                  `json:"height"`
	Version           uint                  `json:"version"`
	CoinBaseValue     uint                  `json:"coinbasevalue"`
	SizeLimit         uint                  `json:"sizelimit"`
	Transactions      []TransactionTemplate `json:"transactions"`
}

type blockTemplate struct {
	Error  string         `json:"error"`
	Result ResultTemplate `json:"result"`
}

type difficultyTemplate struct {
	Error      string  `json:"error"`
	Difficulty float64 `json:"result"`
	ID         string  `json:"id"`
}

// VERY Temporary work-around for GetBlockTemplate() from BP023 ;)
//GetResultTemplate(user, password, host)
//Get and parse data received from Bitcoin client on a getblocktemplate request
func getResultTemplate(user, password, host string) (rtp ResultTemplate, err error) {
	var btp blockTemplate
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

//GetDifficulty function that retrieves current difficulty from bitcoin client
func GetDifficulty(user, password, host string) (difficulty float64, err error) {
	command := "curl -u " + user + ":" + password + ` --data-binary '{"jsonrpc": "1.1", "id":"0", "method": "getdifficulty"}'   -H 'content-type: application/json;' http://` + host + "/"
	var dif difficultyTemplate
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
