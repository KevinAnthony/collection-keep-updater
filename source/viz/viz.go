package viz

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

const (
	baseURL    = "https://www.viz.com"
	sourceName = "Viz"
)

type viz struct {
	settingsHelper
	client http.Client
}

func New(cmd types.ICommand) types.ISource {
	return viz{
		settingsHelper: settingsHelper{},
		client:         ctxu.GetHttpClient(cmd),
	}
}

func (v viz) GetISBNs(ctx context.Context, series types.Series) (types.ISBNBooks, error) {
	settings, err := types.GetSetting[vizSettings](series)
	if err != nil {
		return nil, err
	}

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

	node, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	pages := v.walkSeriesPage(node)
	if settings.MaximumBacklog != nil {
		max := *settings.MaximumBacklog
		if len(pages) > max {
			pages = pages[len(pages)-max:]
		}
	}

	books := types.NewISBNBooks(len(pages))

	for _, page := range pages {
		if settings.Delay != nil {
			time.Sleep(*settings.Delay)
		}

		if book, err := v.getBookFromSeriesPage(ctx, series, page); err != nil {
			return nil, errors.Wrap(err, page)
		} else if book != nil && len(book.ISBN13) > 0 {
			books = append(books, *book)
		}
	}

	for _, blackISBN := range series.ISBNBlacklist {
		index := books.FindIndexByISBN(blackISBN)
		if index >= 0 {
			books = books.RemoveAt(index)
		}
	}

	return books, nil
}

func (v viz) walkSeriesPage(node *html.Node) []string {
	var seriesPages []string

	if node.Type == html.ElementNode && node.Data == "a" && utils.AttrEquals(node.Attr, "role", "presentation") {
		url, found := utils.AttrContains(node.Attr, "href")
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

	node, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	volume := getVolumeFromPath(path)

	return &types.ISBNBook{
		ISBN13: getISBNFromBody(node),
		Title:  fmt.Sprintf("%s: #%s", series.Name, volume),
		Volume: volume,
		Source: sourceName,
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
			return utils.ISBNNormalize(isbn)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		isbn := getISBNFromBody(child)
		if len(isbn) > 0 {
			return isbn
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
