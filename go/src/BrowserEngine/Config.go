package Scrapper

import (
	"fmt"
	"log"
)

// load it into json / sth else so we create it dynamically
const (
	URLS_BUFFER    uint8  = 1<<8 - 1
	URL_PATTERNS   uint8  = 1<<8 - 1
	ELEMENT_BUFFER uint8  = 1<<8 - 1
	BLANK          string = "about:blank"
	DIV            string = "<div>"
	EDIV           string = "</div>"
	IMG            string = "<img>"
	EIMG           string = "</img>"
	H2             string = "<h2>"
	EH2            string = "</h2>"
	SPAN           string = "<span>"
	ESPAN          string = "</span>"
	LI             string = "<li>"
	ELI            string = "</li>"
	OL             string = "<ol>"
	EOL            string = "</ol>"
	A              string = "<a>"
	EA             string = "</a>"
)

type Site struct {
	url      string
	findMap  map[string]string
	visitMap []string
}

type PageInfo struct {
	picture string
	product string
	rating  uint16
	price   uint8
}

func SiteNew(dom map[string]interface{}) Site {
	findMap := make(map[string]string)
	for k, v := range dom {
		if str, ok := v.(string); !ok {
			panic("[FATAL ERROR] Unable to create site variables.")
		} else {
			findMap[k] = str
		}
	}
	visitMap := make([]string, 0, URLS_BUFFER)
	for i := 1; i <= 50; i++ {
		visitMap = append(visitMap, fmt.Sprintf("page-%v.html", i))
	}
	url := findMap["SITE"]
	return Site{
		url:      url,
		findMap:  findMap,
		visitMap: visitMap,
	}
}

/* Initialize the interfaces */
func (self *Site) Url() string {
	if len(self.url) != 0 && IsValidLink(self.url) {
		return self.url
	}
	log.Fatalf("[ERROR] The link seems invalid.\nAccessing: %v.\nStopping.", self.url)
	return ""
}

func (self *Site) GetFindValue(key string) string {
	if len(self.findMap[key]) == 0 {
		fmt.Println("[WARNING] Value under the key is empty.")
	}
	return self.findMap[key]
}

func (self *Site) GetVisitValue() (string, error) {
	if len(self.visitMap[0]) == 0 {
		fmt.Println("[WARNING] Crawling main page directory")
	}

	subdir := self.visitMap[0]
	self.visitMap = self.visitMap[1:]
	return subdir, nil
}
