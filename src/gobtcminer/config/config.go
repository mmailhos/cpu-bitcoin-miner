/*
Author: Mathieu Mailhos
Filename: config.go
Description: Read and parse configuration file to enable a proper connection with the Bitcoin client.
*/
package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type JsonLogger struct {
	Activated bool   `json:"activated"`
	Level     string `json:"level"`
	File      string `json:"file"`
}

type Config struct {
	User     string     `json:"user"`
	Password string     `json:"password"`
	Host     string     `json:"host"`
	Account  string     `json:"account"`
	Threads  int        `json:"threads"`
	Log      JsonLogger `json:"log"`
}

//readconf(filename)
//Read and parse configuration file to enable connection with bitcoin client
func Readconf(filename string) (conf Config) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error:", err)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("Error:", err)
	}
	return
}
