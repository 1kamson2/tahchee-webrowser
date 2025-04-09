package View

import (
	"BrowserEngine/Handlers"
	"BrowserEngine/Utils"
	"errors"
	"net/url"
)

type TAG_STATE uint8

const (
	APPROVE TAG_STATE = 0
	IGNORE  TAG_STATE = 1
)

const (
	LT  byte = 60  // <
	GT  byte = 62  // >
	LSB byte = 91  // [
	RSB byte = 93  // ]
	LRB byte = 40  // (
	RRB byte = 41  // )
	LCB byte = 123 // {
	RCB byte = 125 // }
	N   byte = 10  // \n
)

type Html struct {
	urlInstance url.URL
	body        []byte
}

type HtmlI interface {
	IgnoreStyle(body []byte, bracketMap map[byte]byte) ([]byte, error)
	Lexer()
	Parser()
}

func HtmlNew(url_instance url.URL) (Html, error) {
	if body, err := Handlers.GetRequest(url_instance); err != nil {
		return Html{}, err
	} else {
		return Html{
			urlInstance: url_instance,
			body:        body,
		}, nil
	}
}

func (self *Html) AreTagsCorrect(tagArray []byte, bracketMap map[byte]byte) ([]byte, error) {
	var (
		START_LEXING_FLAG bool = false
		stackSz           int
		tagStack          []byte
	)
	for idx, char := range tagArray {
		if !START_LEXING_FLAG && char == LT {
			START_LEXING_FLAG = true
			tagStack = append(tagStack, char)
			continue
		}
		if _, ok := bracketMap[char]; ok {
			tagStack = append(tagStack, char)
			continue
		}

		stackSz = len(tagStack)
		_, ok := bracketMap[char]
		if stackSz == 0 || (tagStack[stackSz-1] != bracketMap[char] && ok) {
			return []byte{}, errors.New("[ERROR] Invalid CSS.")
		}

		if char == GT {
			tagStack = tagStack[:stackSz-1]
			if stackSz == 1 {
				return self.body[idx+1:], nil
			}
		}
	}
	return []byte{}, errors.New("[ERROR] Out of loop, incorrect CSS")
}

func (self *Html) Lexer(tagStates map[byte]TAG_STATE, bracketMap map[byte]byte) ([]byte, error) {
	var (
		lexed []byte = self.body
		idx   int    = -(1 << 31)
		err   error
		// tagStack    []byte
		// IN_TAG_FLAG bool = false
	)

	/* Sanitize from style if ignore is ON */
	for idx != 0 {
		idx, err = Utils.Grep([]byte("</style>"), lexed)
		if err != nil {
			return []byte{}, err
		}
		if idx != 0 {
			lexed = lexed[idx+9:]
		}
	}

	return lexed, nil
}
