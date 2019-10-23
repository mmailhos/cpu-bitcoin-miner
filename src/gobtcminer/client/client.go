/*Package client for interacting with bitcoin client
Author: Mathieu Mailhos
Filename: client_curl.go
*/
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

// btcClientParamsCaps define a list of capabilities to be ingested in some methods in HTTP requests to the Bitcoin client
type btcClientParamsCaps struct {
	Capabilities []string `json:"capabilities"`
}

// btcClientPayload is a HTTP Payload for the Bitcoin client
type btcClientPayload struct {
	Jsonrpc    string                `json:"jsonrpc"`
	ID         string                `json:"id"`
	Method     string                `json:"method"`
	Parameters []btcClientParamsCaps `json:"params",omitempty`
}

// VERY Temporary work-around for GetBlockTemplate() from BP023 ;)
//GetResultTemplate(user, password, host)
//Get and parse data received from Bitcoin client on a getblocktemplate request
func getResultTemplate(user, password, host string) (rtp ResultTemplate, err error) {
	var btp blockTemplate

	caps := []string{"coinbasetxn", "workid", "coinbase/append"}
	params := btcClientParamsCaps{Capabilities: caps}
	data := btcClientPayload{
		Jsonrpc:    "1.1",
		ID:         "0",
		Method:     "getblocktemplate",
		Parameters: []btcClientParamsCaps{params},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return btp.Result, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://"+host, body)
	if err != nil {
		return btp.Result, err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return btp.Result, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&btp)
	if err != nil {
		return btp.Result, err
	}
	if btp.Error == "" {
		return btp.Result, nil
	}
	return btp.Result, errors.New(btp.Error)
}

//GetDifficulty function that retrieves current difficulty from bitcoin client
func GetDifficulty(user, password, host string) (difficulty float64, err error) {
	var dif difficultyTemplate
	data := btcClientPayload{
		Jsonrpc: "1.1",
		ID:      "0",
		Method:  "getdifficulty",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return dif.Difficulty, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("GET", "http://"+host, body)
	if err != nil {
		return dif.Difficulty, err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dif.Difficulty, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dif)
	if err != nil {
		return dif.Difficulty, err
	}

	if dif.Error == "" {
		return dif.Difficulty, nil
	}
	return dif.Difficulty, errors.New(dif.Error)
}
