package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

var configuration map[string]map[string]string

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

func GetConfigValue(service string, key string) (string, bool) {
	if _, ok := configuration[service]; !ok {
		return "", ok
	}
	v, ok := configuration[service][key]
	return v, ok
}
