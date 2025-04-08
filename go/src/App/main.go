package main

import (
	"BrowserEngine"
)

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
	browser := BrowserEngine.Start()
	browser.Open()

}
