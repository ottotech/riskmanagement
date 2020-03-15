package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type storage int

var StorageType storage
var Logger *log.Logger

const (
	MEMORY storage = iota
	JSON
)

type settings struct {
	Storage string `json:"storage_type"`
	Testing bool   `json:"testing"`
}

func init() {
	var err error
	var s settings

	// get settings from config.json
	goPath := os.Getenv("GOPATH")
	confPath := path.Join(goPath, "/src/github.com/ottotech/riskmanagement/config.json")
	f, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(f, &s)
	if err != nil {
		log.Fatalln(err)
	}
	if s.Storage == "JSON" {
		StorageType = JSON
	}
	if s.Storage == "MEMORY" {
		StorageType = MEMORY
		if s.Testing {
			removeAllMedia()
		}
	}
}
