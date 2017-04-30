package main

import (
	"flag"
	"log"
	"github.com/ti/mdrest"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c",  "", "config file")
	flag.Parse()
	mdj := mdrest.Load(configFile)
	defer mdj.Close()
	if err := mdj.Do(); err != nil {
		log.Panic(err)
	}
}
