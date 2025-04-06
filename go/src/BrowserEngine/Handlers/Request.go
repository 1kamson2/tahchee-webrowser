package Handlers

import (
	"context"
	"github.com/chromedp/chromedp"
	"io"
	"net/http"
	"net/url"
)

func GetRequest(addr url.URL) ([]byte, error) {
	/*
		Parameters:
		addr: Your address that is already the structure URL.
		Function:
		The function makes a request to given site. It doesn't handle any other
		type of actions
	*/
	// TODO: Make headers

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	url := "https://gobyexample.com/enums"
	err := chromedp.Run(ctx, chromedp.Navigate(url), chromedp.WaitVisible("body"), chromedp.TextContent(`body`, &res))
	if err != nil {
		fmt.Println("Error navigating to the website:", err)
		return
	}

	fmt.Println("Website loaded successfully:", res)

	response, err := http.Get(addr.String())
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
