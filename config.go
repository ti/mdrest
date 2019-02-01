package mdrest

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const def_config = `
	{
	  "Watch": false,
	  "BasePath": "",
	  "SrcDir": "./",
	  "OutputType": "json",
	  "SiteMapDeep": 2,
	  "DistDir": "",
	  "NoLogging": false,
	  "NoIndex": false,
	  "NoSiteMap": false,
	  "Server":false,
	  "ServerAddr":":9527"
	}`

type Config struct {
	Watch       bool
	BasePath    string //the base path of you project, default is "", you can use "/" or "/docs/"
	SrcDir      string
	OutputType  string //you can put json or  html is default
	SiteMapDeep int
	DistDir     string
	NoSummary   bool
	NoLogging   bool
	NoIndex     bool
	NoSiteMap   bool
	Server      bool
	ServerAddr  string
}

func LoadConfig(jsonFilePath string) *Config {
	var cfg Config
	err := json.Unmarshal([]byte(def_config), &cfg)
	if err != nil {
		log.Fatal(err)
	}
	if (jsonFilePath != "") {
		bytes, err := ioutil.ReadFile(jsonFilePath)
		if err != nil {
			log.Fatal(err)
		}
		if err = json.Unmarshal(bytes, &cfg); err != nil {
			log.Fatal(err)
		}
	}
	return &cfg
}
