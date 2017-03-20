package mdrest

import (
	"testing"
	"log"
	"encoding/json"
)

func TestReadArticle(t *testing.T) {
	dir := "/Users/leenanxi/go/src/git.tiup.us/go/mdjson/do/content/mm/hwo_to.md"
	ar, err := ReadArticle("",dir)
	if err != nil {
		log.Println(err)
	}
	_ = ar
}

func TestReadArticles(t *testing.T) {
	if 1 == 1 {
		return
	}
	path := "/Users/leenanxi/go/src/git.tiup.us/go/mdjson/do/article/dist/src"
	ars, err := ReadArticles(path,"")
	if err != nil {
		log.Println(err)
	}
	//ars.WriteJsonFiles("/Users/leenanxi/go/src/git.tiup.us/go/mdjson/do/article/dist/dist")

	siteMap := ars.GetSiteMap(2)

	jsb, _ := json.MarshalIndent(siteMap,"","\t")

	log.Println(string(jsb))
}

