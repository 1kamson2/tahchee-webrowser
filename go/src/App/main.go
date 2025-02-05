package main

import (
	"golang.org/Scrapper"
)

func main() {
	/*
	   1. I am interested in seeing the lowest prices of Computer Components,
	   so I would like to scrap every (methaphor) site that sells computer parts
	   and check the lowest prices across the sites. Host it somewhere, make plots.
	   2. Scrapping LinkedIn is hard, but it is worth to see connections between users.
	   3. Scrap wiki and list articles
	   4. Scrap apartment prices
	*/
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
	site := Scrapper.SiteNew(Scrapper.AMAZON_URL)
	run := Scrapper.New(site)
	run.Crawl()

	// fmt.Println(data, err)
}
