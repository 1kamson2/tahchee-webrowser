package main

import (
	"encoding/json"
	"fmt"
	"golang.org/Scrapper"
	"os"
	"path/filepath"
)

func IOCheck(err error) {
	if err != nil {
		fmt.Println(err)
		panic("[FATAL ERROR] Check your IO, failed to read file. ")
	}
}

func JSONConfig(arg string) map[string]interface{} {
	var data map[string]interface{}
	path := filepath.Join(os.Args[3], arg)
	json_bytes, err := os.ReadFile(path)
	IOCheck(err)
	err = json.Unmarshal(json_bytes, &data)
	IOCheck(err)
	return data
}

func main() {
	/*
	   TODO:
	     Use proxies
	     Tor Connections
	     Way to handle gzip encoding
	     Multithreading
	     Multiprocessing
	     Add tokens for parsing
	     Learn about request more.
	     Host it on your server, that allows accessing products.
	     On server should be graph with data
	     (that is the whole idea of this project)
	     Find a way to sort things, based on categories, reviews.
	     In tests prepare some prefetch / predownladed sites.
	*/
	/* How to use?
	   go run main.go "dom.json" "config.json" "path/to/resources/"
	*/
	dom := JSONConfig(os.Args[1])
	agent := JSONConfig(os.Args[2])
	site := Scrapper.SiteNew(dom)
	run := Scrapper.New(site, agent)
	run.Crawl()
}
