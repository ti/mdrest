package main

import (
	"flag"
	"github.com/ti/mdrest"
	"log"
)

var src = flag.String("src", "", "markdown files src dir")
var basePath = flag.String("base", "", "base path for assets url")
var outType = flag.String("out", "html", "use html or json as output")
var configFile = flag.String("c", "", "config file path")
var showTitle = flag.Bool("t", false, "show page title in generated html")

func main() {
	flag.Parse()
	cfg := mdrest.LoadConfig(*configFile)
	if *src != "" {
		cfg.SrcDir = *src
	}
	if *basePath != "" {
		cfg.BasePath = *basePath
	}
	if *outType != "html" {
		cfg.OutputType = *outType
	}
	cfg.ShowPageTitle = *showTitle
	mdr := mdrest.New(cfg)
	defer mdr.Close()
	if err := mdr.Do(); err != nil {
		log.Panic(err)
	}
}
