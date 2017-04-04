package mdrest

import (
	"strings"
	"os"
	"path"
	"path/filepath"
	"encoding/json"
	"log"
	"io/ioutil"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)


const (
	indexName  = "mdrest_index.json"    //full index file, include "title, summery, tags, ohters"
	searchIndexName  = "mdrest_search_index.json" //index just for search, just include  "title, text, tags"
	siteMapName  = "mdrest_sitemap.json"
)

func (this Articles) WriteAllFiles(distDir, fileType string)  {
	this.WriteIndexFile(distDir)
	this.WriteFiles(distDir, fileType)
	this.WriteSiteMapFile(distDir,2)
}

//ReadFiles read .md files and skip ".*" and "_*" file
func ReadFiles(srcDir string) (files []string, err error) {
	walkFunc := func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		filename := path.Base(info.Name())
		if strings.HasPrefix(filename, "_") || strings.HasPrefix(filename, ".")  {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := path.Ext(filename)
		if ext != ".md" {
			return nil
		}
		filename = strings.TrimSuffix(filename, ext)
		files = append(files,fpath)
		return nil
	}
	err = filepath.Walk(srcDir, walkFunc)
	return
}


//WriteJsonFiles
func (this Articles) WriteFiles(distDir string, ftype string)  {
	basePath := ""
	if !strings.HasSuffix(distDir,"/") {
		distDir += "/"
	}
	for _, article := range this {
		dir := path.Join(distDir, basePath)
		arti := *article
		articleLocation :=  arti[KeyLocation].(string)
		relativeDir := path.Dir(articleLocation)
		if relativeDir != "." {
			dir += "/" + relativeDir
		}
		os.MkdirAll(dir, os.ModePerm)
		distFileName := distDir + articleLocation + "." + ftype
		var bytes  []byte
		var err error

		m := minify.New()
		m.AddFunc("text/html", html.Minify)
		htmlContent := []byte(arti[KeyHtml].(string))
		if out,err := m.Bytes("text/html",htmlContent); err == nil {
			htmlContent = out
		}
		if ftype == "json" {
			delete(arti,KeyRawContent)
			bytes, err = json.Marshal(arti)
			arti[KeyHtml] = string(htmlContent)
			if err != nil {
				log.Printf("could not Marshal article %v due to error: %v", distFileName, err)
				continue
			}
		} else if ftype == "html" {
			bytes = htmlContent
		} else {
			bytes = (*article)[KeyRawContent].([]byte)
		}
		writeErr := ioutil.WriteFile(distFileName, bytes, os.ModePerm)
		if writeErr != nil {
			continue
			log.Printf("Could not write file %v due to error: %v", distFileName, writeErr)
		}
	}
}

//WriteSiteMapFile write a file about sitemap
//distDir destination dir
//deep hwo deep to gen site map default is 2
func (this Articles) WriteSiteMapFile(distDir string, deep int) {
	//deep can not be 0
	if deep == 0 {
		deep = 2
	}
	if !strings.HasSuffix(distDir,"/") {
		distDir += "/"
	}
	siteMap := this.GetSiteMap(deep)
	siteMapBytes, _ := json.Marshal(siteMap)
	writeErr := ioutil.WriteFile(distDir + siteMapName, siteMapBytes, os.ModePerm)
	if writeErr != nil {
		log.Printf("Could not write file %v due to error: %v", siteMapName, writeErr)
	}
}


func (this Articles) WriteIndexFile(distDir string) {
	if !strings.HasSuffix(distDir,"/") {
		distDir += "/"
	}
	var indexArticles Articles
	var searchIndexArticles Articles
	for _, article := range this {
		arti :=  make(Article)
		for k, v := range (*article) {
			arti[k] = v
		}
		text := StripHTML((*article)[KeyHtml].(string))
		delete(arti, KeyRawContent)
		delete(arti, KeyHtml)
		arti[KeySummary], arti["truncated"] = StripSummary(text, 80)
		indexArticles = append(indexArticles, &arti)
		//for search
		searchArticle := make(Article)
		searchArticle[KeyText] = text
		searchArticle[KeyTitle] = arti[KeyTitle]
		searchArticle[KeyLocation] = arti[KeyLocation]
		//if it has tags index it
		if tags, ok := arti["tags"]; ok {
			searchArticle["tags"] = tags
		}
		searchIndexArticles =  append(searchIndexArticles,&searchArticle)
	}
	if jsonb, err := json.Marshal(indexArticles); err == nil {
		writeErr := ioutil.WriteFile(distDir + indexName, jsonb, os.ModePerm)
		if writeErr != nil {
			log.Printf("Could not write file %v due to error: %v", indexName, writeErr)
		}
	}
	if jsonb, err := json.Marshal(searchIndexArticles); err == nil {
		writeErr := ioutil.WriteFile(distDir + searchIndexName, jsonb, os.ModePerm)
		if writeErr != nil {
			log.Printf("Could not write file %v due to error: %v", searchIndexName, writeErr)
		}
	}
}


