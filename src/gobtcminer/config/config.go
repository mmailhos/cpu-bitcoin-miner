/*Package config for reading configuration file
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

//JSONLogger structure for writing down logs
type JSONLogger struct {
	Activated bool   `json:"activated"`
	Level     string `json:"level"`
	File      string `json:"file"`
}

//Config simple configuration template from file
type Config struct {
	User     string     `json:"user"`
	Password string     `json:"password"`
	Host     string     `json:"host"`
	Account  string     `json:"account"`
	Log      JSONLogger `json:"log"`
}

//ReadConf reads and parse configuration file to enable connection with bitcoin client
func ReadConf(filename string) (conf Config) {
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
