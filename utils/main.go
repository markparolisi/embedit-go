package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// Store the JSON config as a Go map
var configuration map[string]map[string]string

// Load the credentials.json file if it exists
// Exit with an error if it does not
func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := path.Dir(ex)
	file, _ := os.Open(exPath + "/credentials.json")
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(fmt.Errorf("cannot loading config: %v", err))
	}
}

// Query the configuration map for a service's key
// example imgur.clientID GetConfigValue('imgur', 'clientID')
// Unknown keys will return a blank string and false
func GetConfigValue(service string, key string) (string, bool) {
	if _, ok := configuration[service]; !ok {
		return "", ok
	}
	v, ok := configuration[service][key]
	return v, ok
}


type ErrorMessage struct {
	Code    int
	Message string
}
