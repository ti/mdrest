package main

import (
	"github.com/leenanxi/mdrest"
	"log"
	"flag"
	"io/ioutil"
	"encoding/json"
)


func main() {
	var configFile = flag.String("c",  "config.json", "config file")
	flag.Parse()
	bt, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	var cfg mdrest.Config
	if err = json.Unmarshal(bt, &cfg); err != nil {
		log.Fatal(err)
	}
	mdj := mdrest.New(&cfg)
	defer mdj.Close()
	if err := mdj.Do(); err != nil {
		log.Panic(err)
	}
}
