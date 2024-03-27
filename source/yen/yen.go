package yen

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

const (
	baseURL       = "https://yenpress.com"
	nextQuery     = "next_ord"
	startPosition = "999" // TODO make configurable
	sourceName    = "Yen Press"
)

type yen struct {
	settingsHelper
	client http.Client
}

func New(client http.Client) (types.ISource, error) {
	if client == nil {
		return nil, errors.New("http client is nil")
	}

	return yen{
		settingsHelper: settingsHelper{},
		client:         client,
	}, nil
}

func (y yen) GetISBNs(ctx context.Context, series types.Series) (types.ISBNBooks, error) {
	books, err := y.callGetMore(ctx, series, startPosition)
	if err != nil {
		return nil, err
	}

	sort.Slice(books, func(i, j int) bool {
		iInt, err := strconv.Atoi(books[i].Volume)
		if err != nil {
			return false
		}

		jInt, err := strconv.Atoi(books[j].Volume)
		if err != nil {
			return false
		}

		return iInt < jInt
	})

	for _, blackISBN := range series.ISBNBlacklist {
		index := books.FindIndexByISBN(blackISBN)
		if index >= 0 {
			books = books.RemoveAt(index)
		}
	}

	return books, nil
}

func (y yen) callGetMore(ctx context.Context, series types.Series, next string) (types.ISBNBooks, error) {
	req, err := http.
		NewRequest(y.client).
		Get().
		URL("%s/series/get_more/%s", baseURL, series.ID).
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

	node, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	return y.getBooksFromList(ctx, node, series)
}

func (y yen) getBooksFromList(ctx context.Context, node *html.Node, series types.Series) (types.ISBNBooks, error) {
	var books types.ISBNBooks

	if node.Type == html.ElementNode && node.Data == "a" {
		volumeURL, found := utils.AttrContains(node.Attr, "href")
		if !found {
			return nil, nil
		}
		if len(volumeURL) == 0 {
			dataUrlStr, found := utils.AttrContains(node.Attr, "data-url")
			if !found {
				return nil, nil
			}

			dataUrl, err := url.Parse(dataUrlStr)
			if err != nil {
				return nil, err
			}

			if next := dataUrl.Query().Get(nextQuery); len(next) > 0 {
				return y.callGetMore(ctx, series, next)
			}
		} else {
			book := y.parseURL(volumeURL)
			book.Title = fmt.Sprintf("%s Vol. %s", series.Name, book.Volume)
			books = append(books, book)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		pages, err := y.getBooksFromList(ctx, child, series)
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
		Source: sourceName,
	}
}
