package Handlers

import (
	"fmt"
	"net/url"
)

type TAG_STATE uint8

const (
	APPROVE TAG_STATE = 0
	IGNORE  TAG_STATE = 1
)

const (
	LT byte = 60 // <
	GT byte = 62 // >
	LB byte = 91 // [
	RB byte = 93 // ]
)

type HTML struct {
	URL  url.URL
	Body []byte
	Tags map[byte]TAG_STATE
}

type HTML_I interface {
	ShowHTML() ([]byte, error)
	LexerHTML()
	ParserHTML()
}

func HTMLNew(URL url.URL) (HTML, error) {
	if body, err := GetRequest(URL); err != nil {
		return HTML{}, err
	} else {
		return HTML{
			URL:  URL,
			Body: body,
			Tags: map[byte]TAG_STATE{
				LT: IGNORE,
				GT: IGNORE,
				LB: APPROVE,
				RB: APPROVE,
			},
		}, nil
	}
}

func (self *HTML) LexerHTML() {
	IN_TAG_FLAG := false
	lexed_body := make([]byte, len(self.Body))
	for idx, char := range self.Body {
		switch self.Tags[char] {
		case APPROVE:
			if !IN_TAG_FLAG {
				lexed_body[idx] = char
			}
		case IGNORE:
			// TODO: Add checking if the tags are correct (leetcode)
			IN_TAG_FLAG = !IN_TAG_FLAG
		}
	}
	// is it legal? check if this fucks up the capacity
	self.Body = lexed_body
}

func (self *HTML) ParserHTML() {
	if 1 == 1 {
		fmt.Println("Not implemented")
		return
	}
}

func (self *HTML) ShowHTML() ([]byte, error) {
	return self.Body, nil

}
