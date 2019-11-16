package config

import (
	"encoding/json"
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
	StorageTypeSTR string `json:"storage_type"`
}

func init() {
	var err error
	var s settings

	// get settings from config.json
	goPath := os.Getenv("GOPATH")
	confPath := path.Join(goPath, "/src/github.com/ottotech/riskmanagement/config.json")
	f, err := os.OpenFile(confPath, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&s)
	if err != nil {
		log.Fatalln("Could not decode config.json file; got err: ", err)
	}
	if s.StorageTypeSTR == "JSON" {
		StorageType = JSON
	}
	if s.StorageTypeSTR == "MEMORY" {
		StorageType = MEMORY
		removeAllMedia()
	}
}
