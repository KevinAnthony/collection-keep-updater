package viz

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"
	"golang.org/x/net/html"
	"strconv"
	"strings"
)

const baseURL = "https://www.viz.com"

type viz struct {
	client http.Client
}

func New(client http.Client) types.ISource {
	if client == nil {
		panic("http client is nil")
	}

	return viz{
		client: client,
	}
}

func (v viz) GetISBNs(ctx context.Context, series types.Series) ([]types.ISBNBook, error) {
	req, err := http.
		NewRequest(v.client).
		Get().
		URL("%s/read/manga/%s/all", baseURL, series.ID).
		CreateRequest(ctx)
	if err != nil {
		return nil, err
	}

	body, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	pages := v.walkSeriesPage(node)
	books := make([]types.ISBNBook, 0, len(pages))

	for _, page := range pages {
		if book, err := v.getBookFromSeriesPage(ctx, series, page); err != nil {
			return nil, err
		} else if book != nil {
			books = append(books, *book)
		}

	}

	return books, nil
}

func (v viz) walkSeriesPage(node *html.Node) []string {
	var seriesPages []string

	if node.Type == html.ElementNode && node.Data == "a" && attrEquals(node.Attr, "role", "presentation") {
		url, found := attrContains(node.Attr, "href")
		if found {
			seriesPages = append(seriesPages, url)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		pages := v.walkSeriesPage(child)
		if len(pages) > 0 {
			seriesPages = append(seriesPages, pages...)
		}
	}

	return seriesPages
}

func (v viz) getBookFromSeriesPage(ctx context.Context, series types.Series, path string) (*types.ISBNBook, error) {
	req, err := http.
		NewRequest(v.client).
		Get().
		URL("%s/%s", baseURL, path).
		CreateRequest(ctx)
	if err != nil {
		return nil, err
	}

	body, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	volume := getVolumeFromPath(path)

	return &types.ISBNBook{
		ISBN13:  getISBNFromBody(node),
		Title:   fmt.Sprintf("%s: #%s", series.Name, volume),
		Binding: "",
		Volume:  volume,
	}, nil
}

func getVolumeFromPath(path string) string {
	urlSplit := strings.Split(path, "/")
	if len(urlSplit) < 4 {
		return ""
	}
	stubSplit := strings.Split(urlSplit[3], "-")
	volumeMaybe := stubSplit[len(stubSplit)-1]

	if _, err := strconv.Atoi(volumeMaybe); err == nil { // check if it's a number
		return volumeMaybe
	}

	return ""
}

func getISBNFromBody(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "strong" {
		if isbn := getISBNFromStrong(node); len(isbn) > 0 {
			return strings.ReplaceAll(isbn, "-", "")
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		title := getISBNFromBody(child)
		if len(title) > 0 {
			return title
		}
	}

	return ""
}

func getISBNFromStrong(node *html.Node) string {
	if node.FirstChild == nil || node.NextSibling == nil {
		return ""
	}

	if !strings.EqualFold("ISBN-13", node.FirstChild.Data) {
		return ""
	}

	return node.NextSibling.Data
}

func attrEquals(attr []html.Attribute, key string, value string) bool {
	val, found := attrContains(attr, key)
	if !found {
		return false
	}

	return strings.EqualFold(val, value)
}

func attrContains(attr []html.Attribute, key string) (string, bool) {
	for _, attrKey := range attr {
		if attrKey.Key == key {
			return attrKey.Val, true
		}
	}

	return "", false
}
