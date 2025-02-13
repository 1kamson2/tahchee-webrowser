package Scrapper

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
)

type Stats struct {
	SitesCrawled uint32
	LastCrawled  string
}

type Scrapper struct {
	Client        http.Client
	SiteInstance  Site
	RequestConfig map[string]string
	ScrapperStats Stats
	SaveLocation  string
}

func New(site Site, agent map[string]interface{}) Scrapper {
	agent_cfg := make(map[string]string)
	for k, v := range agent {
		if v == nil {
			agent_cfg[k] = "<nil>"
			continue
		}

		if str, ok := v.(string); !ok {
			panic("[FATAL ERROR] Unable to create agent variables.")
		} else {
			agent_cfg[k] = str
		}
	}
	return Scrapper{http.Client{}, site, agent_cfg, Stats{}, os.Args[3]}
}

func (self *Scrapper) PrepareHeaders(request *http.Request) error {
	skipped := 0
	for k, v := range self.RequestConfig {
		/* Skips any empty value for given key */
		if len(v) == 0 {
			skipped += 1
		} else {
			request.Header.Set(k, v)
		}
	}
	if len(request.Header) != len(self.RequestConfig)-skipped {
		mess := fmt.Sprintf("[ERROR] Something went wrong with preparing headers.\n")
		mess += fmt.Sprintf("Got these:\n")
		for k, v := range request.Header {
			mess += fmt.Sprintf("%v : %v\n", k, v)
		}
		return fmt.Errorf(mess)
	}
	return nil
}

func (self Scrapper) GetRequest(subUrl string) (string, error) {
	request, err := http.NewRequest("GET", subUrl, nil)
	if err != nil {
		return "", fmt.Errorf("[GET Request Failed]: %s\n", err)
	}
	err = self.PrepareHeaders(request)
	if err != nil {
		return "", err
	}
	response, err := self.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("[GET Failed]: %v\n", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[GET Failed] Status code: %d\n", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("[GET Failed] Body: %s\n", err)
	}

	return string(body), nil
}

func IsValidLink(subUrl string) bool {
	/* Regex in json, check if it will fuck up the entire program */
	if pattern, err := regexp.Compile(`^(?:https?://)?(?:[^/.\s]+\.)*books.toscrape\.com(?:/[^/\s]+)*/?$`); err == nil {
		return pattern.MatchString(subUrl)
	}
	return false
}

func Grep(line, pat string) int {
	if len(pat) == 0 || len(line) == 0 {
		return -1
	}
	shiftTable := [ELEMENT_BUFFER]int{}

	for i := 0; i < int(ELEMENT_BUFFER); i++ {
		shiftTable[i] = -1
	}
	for i := 0; i < len(pat); i++ {
		shiftTable[pat[i]] = i
	}
	m, n := len(pat), len(line)
	var s int = 0
	for s <= n-m {
		var j int = m - 1
		for j >= 0 && pat[j] == line[s+j] {
			j--
		}
		if j < 0 {
			return s
		} else {
			s += max(1, j-shiftTable[line[s+j]])
		}
	}
	return -1
}

func (self *Scrapper) Lexer(body string) ([]string, error) {
	/*
			   Lexer will analyze the HTML Document structure. It will recognize tokens
			   such as:
			     > <a ... /a>
			     > <div ... /div>
			   The scanning will be using Boyer-Moore algorithm.
			   The meaning of masks (bit shifts):
		      0: denotes the entry point of the page id="a-page"
			    1-2: denotes the entry point - class="sg-col-inner", we skip it twice:
			    3: denotes the starting point of lexing the product, <div> is added to
			       tokens (current stack)
			    4: denotes that we are looking for a picture of the product.
			    5: denotes that we are looking for a product's name
			    6: denotes that we are looking for a product's price
			    7: denotes that we are looking for a product's rating
			    default: we reset to the second bit
	*/
	/* Add Update function to handle those things */

	var mask uint8 = 1
	var toks []string = make([]string, 0, ELEMENT_BUFFER)
	var isAll int = 0
	bodysplt := strings.Split(body, "\n")
	for _, line := range bodysplt {
		/* Sanitizing line from whitespaces */
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		switch mask {
		case 1 << 0:
			// skip
			if index := Grep(line, self.SiteInstance.GetFindValue("SITE_ENTRY")); index != -1 {
				mask = mask << 1
			}
		case 1 << 1:
			// skip
			if index := Grep(line, self.SiteInstance.GetFindValue("PRODUCTS_ENTRY")); index != -1 {
				mask = mask << 1
			}
		case 1 << 2:
			// skip
			if index := Grep(line, self.SiteInstance.GetFindValue("PRODUCT_ENTRY")); index != -1 {
				toks = append(toks, DIV)
				mask = mask << 1
			}
		case 1 << 3:
			if index := Grep(line, self.SiteInstance.GetFindValue("PRODUCT_IMG")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN)
				isAll++
				mask = mask << 1
			}
		case 1 << 4:
			if index := Grep(line, self.SiteInstance.GetFindValue("PRODUCT_NAME")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN)
				isAll++
				mask = mask << 1
			}
		case 1 << 5:
			if index := Grep(line, self.SiteInstance.GetFindValue("PRODUCT_RATING")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN)
				isAll++
				mask = mask << 1
			}
		case 1 << 6:
			if index := Grep(line, self.SiteInstance.GetFindValue("PRODUCT_PRICE")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN, EDIV)
				isAll++
				mask = 1 << 2
			} else if index := Grep(line, self.SiteInstance.GetFindValue("NEXT_PAGE")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN, EDIV)
				isAll++
			}
		default:
			return []string{}, fmt.Errorf("[ERROR] Invalid options, skipping this page.")
		}
	}
	if len(toks) == 0 {
		return []string{}, fmt.Errorf("[ERROR] Didn't parse everything.")
	}

	return toks, nil
}

func (self *Scrapper) ParserSanitize(tok, pat string) []string {
	repat := regexp.MustCompile(pat)
	matpat := repat.FindStringSubmatch(tok)
	return matpat
}

func (self *Scrapper) Parser(lexed []string) (string, error) {
	/*
	   After Lexer work, parser will recognize tokens and work with them
	   to extract data.
	*/

	validBToks := []string{DIV, IMG, H2, SPAN, LI, OL, A}
	validEToks := []string{EDIV, EIMG, EH2, ESPAN, ELI, EOL, EA}
	//validMToks := map[string]string{EDIV: DIV, EIMG: IMG, EH2: H2, ESPAN: SPAN, ELI: LI, EOL: OL, EA: A}
	var lastTok string = lexed[0]
	if !slices.Contains(validBToks, lastTok) {
		return "",
			fmt.Errorf("[ERROR]: This token doesn't start the document: %v\n", lastTok)
	}

	var htmlBody string = fmt.Sprintf("%v\n",
		`<!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Book Store</title>
      <link rel="stylesheet" href="styles.css">    
    </head>
    <body>`)
	var mask uint8 = 1
	for _, tok := range lexed {
		/* Check if we start with the correct token */
		if slices.Contains(validBToks, tok) || slices.Contains(validEToks, tok) {
			htmlBody += fmt.Sprintf("%v\n", tok)
			lastTok = tok
		} else {
			/*
			   Extract the data instead.
			   The lexer adds the whole line, where we encountered the given keyword.
			   Therefore we need to sanitize the whole line. Using the regular
			   expressions, that user provided
			*/
			switch mask {
			case 1 << 0:
				pat := self.SiteInstance.GetFindValue("PRODUCT_IMG_RGX")
				domain := self.SiteInstance.GetFindValue("DOMAIN")
				curr_site := self.SiteInstance.Url()
				elements := self.ParserSanitize(tok, pat)
				htmlEl := fmt.Sprintf("<a href=\"%v%v\"><img alt=\"%v\" src=\"%v%v\"></a>",
					curr_site, elements[1], elements[3], domain, elements[2][2:])
				htmlBody += fmt.Sprintf("%v\n", htmlEl)
				mask = mask << 1
			case 1 << 1:
				pat := self.SiteInstance.GetFindValue("PRODUCT_NAME_RGX")
				elements := self.ParserSanitize(tok, pat)
				htmlEl := fmt.Sprintf("%v %v %v", H2, elements[2], EH2)
				htmlBody += fmt.Sprintf("%v\n", htmlEl)
				mask = mask << 1
			case 1 << 2:
				pat := self.SiteInstance.GetFindValue("PRODUCT_RATING_RGX")
				elements := self.ParserSanitize(tok, pat)
				htmlEl := fmt.Sprintf("Rating: %v", elements[1])
				htmlBody += fmt.Sprintf("%v\n", htmlEl)
				mask = mask << 1
			case 1 << 3:
				pat := self.SiteInstance.GetFindValue("PRODUCT_PRICE_RGX")
				elements := self.ParserSanitize(tok, pat)
				htmlEl := fmt.Sprintf("Price: $%v", elements[1])
				htmlBody += fmt.Sprintf("%v\n", htmlEl)
				mask = 1
			}
		}
	}
	htmlBody += "</body></html>"
	return htmlBody, nil
}

func (self *Scrapper) Update(recentUrl, body string) (string, error) {
	/*
		Update() will add records of the current crawled site. Also we will update
		state of our program. The results of Lexing, Parsing will be passed here.
	*/
	self.ScrapperStats.SitesCrawled++
	self.ScrapperStats.LastCrawled = recentUrl
	filename := uuid.New().String()
	resources := self.SiteInstance.GetFindValue("RESOURCE_DIRECTORY")
	site, err := os.Create(fmt.Sprintf("%v%v.html", resources,
		strings.Replace(filename, "-", "", -1)))
	if err != nil {
		return "", fmt.Errorf("[IO ERROR] Failed to save the page.")
	}
	site.WriteString(body)
	return filename, nil
}

func (self *Scrapper) InfoFetch() {
	fmt.Println("===========================")
	fmt.Printf("CLIENT INFO: %v\n", self.Client)
	fmt.Printf("URL: %v\n", self.SiteInstance.Url)
	fmt.Println("REQUEST CONFIG:")
	for k, v := range self.RequestConfig {
		fmt.Printf("  %v : %v\n", k, v)
	}
	fmt.Printf("Save Location: %v\n", self.SaveLocation)
	fmt.Println("===========================")
}

func (self *Scrapper) Crawl() {
	self.InfoFetch()
	var domainUrl string = self.SiteInstance.Url()
	var newUrl string
	for {
		if subdir, err := self.SiteInstance.GetVisitValue(); err != nil {
			fmt.Printf("[WARNING]: %v is invalid.\n", newUrl)
			continue
		} else {
			newUrl = domainUrl + subdir
			fmt.Printf("[Crawl()] Currently visiting: %v\n", newUrl)
		}

		if !IsValidLink(newUrl) {
			fmt.Printf("[WARNING]: %v is invalid.\n", newUrl)
			continue
		}

		body, err := self.GetRequest(newUrl)
		if err != nil {
			fmt.Printf("[WARNING] Crawl(): %v\nContinue to crawl...\n", err)
			continue
		}

		lexed, err := self.Lexer(body)
		if err != nil {
			fmt.Printf("[WARNING] Crawl(): Error encountered while lexing.\n%v", err)
		}

		parsed, err := self.Parser(lexed)
		if err != nil {
			fmt.Printf("[WARNING] Crawl(): Error encountered while parsing.\n%v", err)
			panic("TODO")
		}

		filename, err := self.Update(newUrl, parsed)
		if err != nil {
			panic("[CRITICAL ERROR] Failed to save the page, might be IO problem.")
		} else {
			fmt.Printf("[SUCCESS] Parsed site saved in: %v\n", filename)
		}

		if self.ScrapperStats.SitesCrawled == uint32(URLS_BUFFER) {
			panic("Nothing else to do")
		}
	}
}
