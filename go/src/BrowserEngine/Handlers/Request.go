package Handlers

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
)

func GetRequest(addr url.URL) ([]byte, error) {
	/*
		Parameters:
		addr: Your address that is already the structure URL.
		Function:
		The function makes a request to given site. It doesn't handle any other
		type of actions.
	*/
	ctx_alloc, cancel_alloc := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
		)...,
	)
	defer cancel_alloc()

	ctx, cancel := context.WithTimeout(ctx_alloc, 10*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		chromedp.Navigate(addr.String()),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.InnerHTML("html", &html),
	)
	if err != nil {
		return []byte{}, err
	}

	return []byte(html), nil
}

func GetRequestSimple(addr url.URL) ([]byte, error) {
	/*
		Parameters:
		addr: Your address that is already the structure URL.
		Function:
		The function makes a request to given site. This should be used as
		a last resort if the main function for getting requests fails.
	*/
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
