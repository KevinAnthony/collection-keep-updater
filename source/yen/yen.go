package yen

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/source"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"
	"golang.org/x/net/html"
)

const (
	baseURL       = "https://yenpress.com"
	nextQuery     = "next_ord"
	startPosition = "999" // TODO make configurable
)

type yen struct {
	client http.Client
}

func New(client http.Client) types.ISource {
	if client == nil {
		panic("http client is nil")
	}

	return yen{
		client: client,
	}
}

func (y yen) GetISBNs(ctx context.Context, series types.Series) (types.ISBNBooks, error) {
	return y.callGetMore(ctx, series.ID, startPosition)
}

func (y yen) callGetMore(ctx context.Context, id, next string) (types.ISBNBooks, error) {
	req, err := http.
		NewRequest(y.client).
		Get().
		URL("%s/series/get_more/%s", baseURL, id).
		Header("x-requested-with", "XMLHttpRequest").
		Query(nextQuery, next).
		CreateRequest(ctx)
	if err != nil {
		return nil, err
	}

	body, err := y.client.Do(req)
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return y.getBooksFromList(ctx, node, id)
}

func (y yen) getBooksFromList(ctx context.Context, node *html.Node, id string) (types.ISBNBooks, error) {
	var books types.ISBNBooks

	if node.Type == html.ElementNode && node.Data == "a" {
		volumeURL, found := source.AttrContains(node.Attr, "href")
		if !found {
			return nil, nil
		}
		if len(volumeURL) == 0 {
			dataUrlStr, found := source.AttrContains(node.Attr, "data-url")
			if !found {
				return nil, nil
			}

			dataUrl, err := url.Parse(dataUrlStr)
			if err != nil {
				return nil, err
			}

			if next := dataUrl.Query().Get(nextQuery); len(next) > 0 {
				return y.callGetMore(ctx, id, next)
			}
		} else {
			books = append(books, y.parseURL(volumeURL))
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		pages, err := y.getBooksFromList(ctx, child, id)
		if err != nil {
			return nil, err
		}

		if len(pages) > 0 {
			books = append(books, pages...)
		}
	}

	return books, nil
}

func (y yen) parseURL(url string) types.ISBNBook {

	fmt.Println(url)
	slugs := strings.Split(url, "/")
	slugSplit := strings.Split(slugs[len(slugs)-1], "-")
	isbn13 := slugSplit[0]
	vol := ""

	for i := range slugSplit {
		if slugSplit[i] != "vol" {
			continue
		}

		if len(slugSplit)-1 == i {
			break // if we at the end, we did not find a volume
		}

		vol = slugSplit[i+1]

		break
	}

	return types.ISBNBook{
		ISBN13: isbn13,
		Volume: vol,
	}
}
