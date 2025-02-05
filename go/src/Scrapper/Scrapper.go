package Scrapper

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"slices"
	//"log"
	"net/http"
	"regexp"
	"strings"
)

type Stats struct {
	SitesCrawled uint32
	LastCrawled  string
}

type Scrapper struct {
	Client         http.Client
	ScrapperConfig SiteInterface
	RequestConfig  map[string]string
	ScrapperStats  Stats
	SaveLocation   string
}

func New(cfg SiteInterface) Scrapper {
	req_config := map[string]string{
		"Host":            "",
		"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:134.0) Gecko/20100101 Firefox/134.0",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		"Connection":      "keep-alive",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
		// No use of encoding for now.
		//"Accept-Encoding": "gzip, deflate",
	}
	return Scrapper{http.Client{}, cfg, req_config, Stats{}, ""}
}

func (self *Scrapper) PrepareHeaders(request *http.Request) error {
	skipped := 0
	for k, v := range self.RequestConfig {
		/* Skips any empty value for given key */
		if len(v) == 0 {
			skipped += 1
		} else {
			request.Header.Add(k, v)
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
	fmt.Printf("%v", request)
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
	if pattern, err := regexp.Compile(`^(?:https?://)?(?:[^/.\s]+\.)*amazon\.com(?:/[^/\s]+)*/?$`); err == nil {
		return pattern.MatchString(subUrl)
	}
	return false
}

func Grep(line, pat string) int {
	if len(pat) == 0 || len(line) == 0 {
		return -1
	}
	shiftTable := [MAX_ELEMENT_BUFFER]int{}

	for i := 0; i < int(MAX_ELEMENT_BUFFER); i++ {
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
	var toks []string = make([]string, 0, MAX_ELEMENT_BUFFER)
	var isAll int = 0
	bodysplt := strings.Split(body, "\n")
	for _, line := range bodysplt {
		/* Sanitizing line from whitespaces */
		line = strings.TrimSpace(line)
		switch mask {
		case 1 << 0:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("PAGE_ENTRY")); index != -1 {
				mask = mask << 1
			}
		case 1 << 1:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("PRODUCT_ENTRY")); index != -1 {
				mask = mask << 1
			}
		case 1 << 2:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("SKIP")); index != -1 {
				mask = mask << 1
			}
		case 1 << 3:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("PRODUCT_IN")); index != -1 {
				toks = append(toks, DIV, line, EDIV)
				isAll++
				mask = mask << 1
			}
		case 1 << 4:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("PRODUCT_PICTURE")); index != -1 {
				toks = append(toks, IMG, line, EIMG)
				isAll++
				mask = mask << 1
			}
		case 1 << 5:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("PRODUCT_NAME")); index != -1 {
				toks = append(toks, H2, line, EH2)
				isAll++
				mask = mask << 1
			}
		case 1 << 6:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("PRODUCT_PRICE")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN)
				isAll++
				mask = mask << 1
			} else if index := Grep(line, self.ScrapperConfig.GetFindValue("PRODUCT_PRICE_USED")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN)
				isAll++
				mask = mask << 1
			}
		case 1 << 7:
			if index := Grep(line, self.ScrapperConfig.GetFindValue("PRODUCT_RATING")); index != -1 {
				toks = append(toks, SPAN, line, ESPAN)
				isAll++
			}
		default:
			return []string{}, fmt.Errorf("[ERROR] Invalid options, skipping this page.")
		}
	}

	if len(toks) != 3*isAll {
		return []string{}, fmt.Errorf("[ERROR] Didn't parse everything.")
	}

	return toks, nil
}

func (self *Scrapper) Parser(lexed []string) (string, error) {
	/*
	   After Lexer work, parser will recognize tokens and work with them
	   to extract data.
	*/

	validBToks := []string{DIV, IMG, H2, SPAN}
	validEToks := []string{EDIV, EIMG, EH2, ESPAN}
	validMToks := map[string]string{EDIV: DIV, EIMG: IMG, EH2: H2, ESPAN: SPAN}

	var lastTok string = lexed[0]
	if !slices.Contains(validBToks, lastTok) {
		return "",
			fmt.Errorf("[ERROR]: This token doesn't start the document: %v\n", lastTok)
	}

	lexed = lexed[1:]
	var htmlBody string = fmt.Sprintf("%v\n", lastTok)
	for _, tok := range lexed {
		if slices.Contains(validBToks, tok) {
			htmlBody += fmt.Sprintf("%v\n", tok)
			lastTok = tok
		} else if slices.Contains(validEToks, tok) {
			if lastTok == validMToks[tok] {
				htmlBody += fmt.Sprintf("%v\n", tok)
			} else {
				return "", fmt.Errorf("[ERROR]: This token doesn't match the map:\n %v != $v", lastTok, validMToks[tok])
			}
		} else {
			// TODO: Find a function that indents or make one
			// TODO: Actually you should extract here data.
			if len(lastTok) < len(tok) {
				start := len(lastTok)
				end := len(tok) - start - 1
				htmlBody += fmt.Sprintf("  %v", tok[start:end])
			} else {
				return "", fmt.Errorf("[ERROR] Tokens don't match:\n%v\n%v\n", tok, lastTok)
			}
		}
	}
	return htmlBody, nil
}

func (self *Scrapper) Update(recentUrl, body string) {
	/*
		Update() will add records of the current crawled site. Also we will update
		state of our program. The results of Lexing, Parsing will be passed here.
	*/
	self.ScrapperStats.SitesCrawled++
	self.ScrapperStats.LastCrawled = recentUrl
	site, _ := os.Create(fmt.Sprintf("%v%v.html", RESOURCE_DIR,
		strings.Replace(uuid.New().String(), "-", "", -1)))
	site.WriteString(body)
}

func (self *Scrapper) Crawl() {
	var domainUrl string = self.ScrapperConfig.Url()
	var newUrl string
	for {
		if subdir, err := self.ScrapperConfig.GetVisitValue(); err != nil {
			fmt.Printf("[WARNING]: %v is invalid.\n", newUrl)
			continue
		} else {
			newUrl = domainUrl + subdir
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
		}

		self.Update(newUrl, parsed)

	}
}
