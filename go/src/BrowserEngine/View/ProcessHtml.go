package View

import (
	"BrowserEngine/Handlers"
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
func (self *Html) Lexer(
	tagStates map[byte]TAG_STATE,
	bracketMap map[byte]byte) ([]byte, error) {
	lexedBody := make([]byte, 0, len(self.body))
	tagStack := make([]byte, 0)
	var (
		lexedIdx        int  = 0
		tagStackSize    int  = 0
		IN_TAG_FLAG     bool = false
		CHECK_TAGS_FLAG bool = true
	)
	for _, char := range self.body {
		state, ok := tagStates[char]
		if ok && state == IGNORE && CHECK_TAGS_FLAG {
			tagStack = append(tagStack, char)
			tagStackSize++
			IN_TAG_FLAG = true
			continue
		} else if ok && state == APPROVE && CHECK_TAGS_FLAG {
			tagStack = append(tagStack, char)
			tagStackSize++
		}
		if CHECK_TAGS_FLAG && ((tagStackSize == 0 && !IN_TAG_FLAG) || tagStack[tagStackSize-1] != bracketMap[char]) {
			return []byte{}, errors.New("[ERROR] The HTML Document is not correct")
		}

		if tagStackSize == 0 && IN_TAG_FLAG && CHECK_TAGS_FLAG {
			IN_TAG_FLAG = false
			CHECK_TAGS_FLAG = false

		}

		if !IN_TAG_FLAG {
			if lexedIdx > 0 && lexedIdx%64 == 0 {
				lexedBody[lexedIdx] = N
				lexedIdx++
				lexedBody[lexedIdx] = char
				lexedIdx++
			} else {
				lexedBody[lexedIdx] = char
				lexedIdx++
			}
		}

		if CHECK_TAGS_FLAG {
			tagStack = tagStack[:tagStackSize-1]
			tagStackSize--
		}
	}
	self.body = lexedBody
	return self.body, nil
}
