package Scrapper

import (
	"fmt"
	"log"
)

const (
	MAX_URL_BUFFER     uint8  = 1<<8 - 1
	MAX_URL_PATTERNS   uint8  = 1<<8 - 1
	MAX_ELEMENT_BUFFER uint8  = 1<<8 - 1
	AMAZON_URL         string = "https://www.amazon.com/"
	LINKEDIN_URL       string = "https://linkedin.com/"
	PAGE_ENTRY         string = "<div id=\"a-page\">"
	PRODUCT_ENTRY      string = "<div class=\"sg-col-inner\">"
	SKIP               string = "<div class=\"sg-col-inner\">"
	PRODUCT_IN         string = "<div class=\"sg-col-inner\">"
	PRODUCT_PICTURE    string = "<img class=\"s-image\" src="
	PRODUCT_NAME       string = "<h2 class=\"a-size-base-plus a-spacing-none a-color-base a-text-normal\" aria-label="
	PRODUCT_PRICE      string = "<span class\"a-offscreen\">"
	PRODUCT_PRICE_USED string = "<span class=\"a-color-base\">"
	PRODUCT_RATING     string = "<span class=\"a-size-base s-underline-text\" aria-hidden=\"true\""
	DIV                string = "<div>"
	EDIV               string = "</div>"
	IMG                string = "<img>"
	EIMG               string = "</img>"
	H2                 string = "<h2>"
	EH2                string = "</h2>"
	SPAN               string = "<span>"
	ESPAN              string = "</span>"
	RESOURCE_DIR       string = "/home/kums0nd/Dev/scrapper/go/resources/"
)

type SiteInterface interface {
	Url() string
	GetFindValue(key string) string
	GetVisitValue() (string, error)
	// return members function
}

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

func SiteNew(url string) Site {
	/*
	 A starting link is a first page in category. An example:
	 https://www.amazon.com/s?i=computers-intl-ship&bbn=16225007011&rh=n%3A16
	 225007011%2Cn%3A193870011&qid=1738360354&xpid=_vIRAGNHTpOQ-&ref=sr_pg_1
	 1. An entry point seems to be line starting with: <div id="a-page"> </div>,
	 everything inside this div is a related to products, that you will
	 be scraping.
	 2. Each product is placed in: <div class="sg-col-inner"></div>, but the
	 first one we encounter is actually a container for every product on
	 the current site. Probably some sort of flag might be needed.
	 3. If we wanted an picture of the product, we could access it by:
	  <img class="s-image" src="https://m.media-amazon.com/images/I/71aHvYUgX1L._AC_UL320_.jpg"
	 4. To get the product's name we get to the line:
	  <h2 aria-label="AMD RYZEN 7 9800X3D 8-Core, 16-Thread Desktop Processor"
	  class="a-size-base-plus a-spacing-none a-color-base a-text-normal">
	  <span>AMD RYZEN 7 9800X3D 8-Core, 16-Thread Desktop Processor</span></h2>
	 We extract data from aria-label or from span.
	 5. Prices might be displayed in two different ways:
	    5.1) <span class="a-price-whole" ... /> and
	    in class "a-price-fraction"
	    5.2) <span class="a-color-base>, if the previous one is not present.
	    We are interested only in one of them.
	 6. For ratings we access:  <a aria-label="318 ratings" class="a-link-normal s-underline-text s-underline-link-text s-link-style">
	    We extract ratings.
	*/

	// add handling for other stuff

	findMap := map[string]string{
		"PAGE_ENTRY":         PAGE_ENTRY,
		"PRODUCT_ENTRY":      PRODUCT_ENTRY,
		"PRODUCT_IN":         PRODUCT_IN,
		"PRODUCT_PICTURE":    PRODUCT_PICTURE,
		"PRODUCT_NAME":       PRODUCT_NAME,
		"PRODUCT_RATING":     PRODUCT_RATING,
		"PRODUCT_PRICE":      PRODUCT_PRICE,
		"PRODUCT_PRICE_USED": PRODUCT_PRICE_USED,
	}

	visitMap := make([]string, 0, MAX_URL_BUFFER)
	visitMap = append(visitMap,
		"s?i=specialty-aps&bbn=16225007011&rh=n%3A16225007011%2Cn%3A193870011&ref=nav_em__nav_desktop_sa_intl_computer_components_0_2_7_3")
	// this is how we add
	// div = append("<div class=", "<div class='A'")
	return Site{
		url:      url,
		findMap:  findMap,
		visitMap: visitMap,
	}
}

/* Initialize the interfaces */
func (self Site) Url() string {
	if len(self.url) != 0 && IsValidLink(self.url) {
		return self.url
	}
	log.Fatalf("[ERROR] The link seems invalid.\nAccessing: %v.\nStopping.", self.url)
	return ""
}

func (self Site) GetFindValue(key string) string {
	return self.findMap[key]
}

func (self Site) GetVisitValue() (string, error) {
	if len(self.visitMap[0]) == 0 {
		return "", fmt.Errorf("[WARNING] This entry in the buffer is empty.")
	}

	subdir := self.visitMap[0]
	self.visitMap = self.visitMap[1:]
	return subdir, nil
}
