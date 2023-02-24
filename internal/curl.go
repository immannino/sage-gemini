package internal

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/yterajima/go-sitemap"
)

func FetchSitemap(url string) (*sitemap.Sitemap, error) {
	sm, err := sitemap.Get(url, nil)
	if err != nil {
		return nil, err
	}

	return &sm, nil
}

func FetchHTML(ctx context.Context, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > http.StatusUnsupportedMediaType {
		return "", errors.New("Not successful response from server")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
