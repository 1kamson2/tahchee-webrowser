package BrowserEngine

import (
	"BrowserEngine/View"
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

var (
	/* Those are all constants, shouldn't be changed */
	// TODO: This ain't working
	VALID_PROTOCOL = regexp.MustCompile(`^https?:\/\/`)
	VALID_WWW      = regexp.MustCompile(`(www\.)?`)
	VALID_DOMAIN   = regexp.MustCompile(`[a-zA-Z0-9-]+\.[a-zA-Z0-9-]+`)
	VALID_PORT     = regexp.MustCompile(`(:\d{1,5})?`)
	VALID_TLD      = regexp.MustCompile(`\.[a-zA-Z]{2,}`)
	VALID_PATH     = regexp.MustCompile(`(\/.*)?$`)
)

type Browser struct {
	tabs       []View.Html
	bracketMap map[byte]byte
	tagStates  map[byte]View.TAG_STATE
}

type BrowserI interface {
	Lexer()
	Parser()
	Show() ([]byte, error)
	PreprocessUrl(url_s string) (url.URL, error)
	Start()
}

func Start() Browser {
	return Browser{
		tabs: []View.Html{},
		bracketMap: map[byte]byte{
			View.RRB: View.LRB,
			View.RSB: View.LSB,
			View.RCB: View.LCB,
		},
		tagStates: map[byte]View.TAG_STATE{
			View.LT:  View.APPROVE,
			View.GT:  View.APPROVE,
			View.LSB: View.APPROVE,
			View.RSB: View.APPROVE,
			View.LRB: View.APPROVE,
			View.RRB: View.APPROVE,
			View.LCB: View.IGNORE,
			View.RCB: View.IGNORE,
		},
	}
}

func (self *Browser) Parser() {
	if true {
		fmt.Println("Not implemented")
		return
	}
}

func (self *Browser) Show() ([]byte, error) {
	return []byte{}, nil

}

func (self *Browser) PreprocessUrl(url_s string) (url.URL, error) {
	if !VALID_PROTOCOL.MatchString(url_s) {
		/* The beginning of the link is not https://,
		so we check if it is the www. */
		if VALID_WWW.MatchString(url_s) {
			url_s = fmt.Sprintf("https://%v", url_s)
		} else {
			/* The url is broken, go back */
			return url.URL{}, errors.New(
				fmt.Sprintf("[ERROR] The url is broken. Got '%v'.\n", url_s))
		}
	}

	if !VALID_DOMAIN.MatchString(url_s) {
		return url.URL{}, errors.New(
			fmt.Sprintf("[ERROR] Invalid domain. Got '%v'.\n", url_s))
	}

	if !VALID_TLD.MatchString(url_s) {
		return url.URL{}, errors.New(
			fmt.Sprintf("[ERROR] Invalid TLD. Got '%v'.\n", url_s))
	}

	if !VALID_PATH.MatchString(url_s) {
		return url.URL{}, errors.New(
			fmt.Sprintf("[ERROR] Trying to access invalid resources. Got '%v'.\n", url_s))
	}

	url_instance, err := url.Parse(url_s)
	if err != nil {
		return url.URL{}, err
	}
	return *url_instance, nil
}

func (self *Browser) Open() {
	var (
		SEARCH_BAR_FLAG bool = true
		DO_ACTIONS_FLAG bool = false
	)
	for {
		if SEARCH_BAR_FLAG {
			url_s := View.SearchBar()
			url_instance, _ := self.PreprocessUrl(url_s)
			html, err := View.HtmlNew(url_instance)

			if err != nil {
				fmt.Println(err)
				continue
			} else {
				SEARCH_BAR_FLAG = false
			}
			body, err := html.Lexer(self.tagStates, self.bracketMap)
			if err != nil {
				fmt.Println(err)
			}
			for _, char := range body {
				fmt.Printf("%v", string(char))
			}
			DO_ACTIONS_FLAG = true
		}
		if DO_ACTIONS_FLAG {
			fmt.Println("Hello here action")
			SEARCH_BAR_FLAG = true
		}
	}
}
