package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration config of the application
type Configuration struct {
	Appname  string   `json:"appname"`
	Address  string   `json:"address"`
	Port     int      `json:"port"`
	Static   string   `json:"static"`
	Acme     bool     `json:"acme"`
	Acmehost []string `json:"acmehost"`
	DirCache string   `json:"dirCache"`
	Crt      string   `json:"crt,omitempty"`
	Key      string   `json:"key,omitempty"`
}

var (
	conf    *Configuration
	csrfkey string
)

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		Danger("Cannot open config.json file", err)
	}

	decoder := json.NewDecoder(file)
	conf = &Configuration{}
	err = decoder.Decode(conf)
	if err != nil {
		Danger("Cannot get configuration from file", err)
	}
}

// Env gets configuration
func Env(reload bool) *Configuration {
	if reload {
		loadConfig()
	}

	return conf
}

func setCsrf() {
	csrfkey = CreateUUID()

	fmt.Println(csrfkey)
}

func GetKeyCSRF() string {
	return csrfkey
}
