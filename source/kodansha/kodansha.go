package kodansha

import (
	"context"
	"strconv"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

const (
	baseURL    = "https://kodansha.us"
	seriesSlug = "/series/"
	sourceName = "Kodansha"
)

type kodansha struct {
	settingsHelper
	client http.Client
}

func New(client http.Client) (types.ISource, error) {
	if client == nil {
		return nil, errors.New("http client is nil")
	}

	return kodansha{
		settingsHelper: settingsHelper{},
		client:         client,
	}, nil
}

func (k kodansha) GetISBNs(ctx context.Context, series types.Series) (types.ISBNBooks, error) {
	req, err := http.
		NewRequest(k.client).
		Get().
		URL("%s%s%s", baseURL, seriesSlug, series.ID).
		CreateRequest(ctx)
	if err != nil {
		return nil, err
	}

	body, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	pages := k.walkSeriesPage(node)

	books := types.NewISBNBooks(len(pages))

	for _, page := range pages {
		if book, err := k.getBookFromSeriesPage(ctx, page); err != nil {
			return nil, errors.Wrap(err, page)
		} else if book != nil {
			books = append(books, *book)
		}
	}
	return books, nil
}

func (k kodansha) walkSeriesPage(node *html.Node) []string {
	var seriesPages []string

	if node.Type == html.ElementNode && node.Data == "span" && utils.AttrEquals(node.Attr, "class", "details-text") {
		if detail := k.extractDetails(node.FirstChild); len(detail) > 0 {
			seriesPages = append(seriesPages, detail)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		pages := k.walkSeriesPage(child)
		if len(pages) > 0 {
			seriesPages = append(seriesPages, pages...)
		}
	}

	return seriesPages
}

func (k kodansha) extractDetails(node *html.Node) string {
	for child := node; child != nil; child = child.NextSibling {
		if child.Data == "a" {
			url, found := utils.AttrContains(child.Attr, "href")
			if found {
				return url
			}
		}
	}

	return ""
}

func (k kodansha) getBookFromSeriesPage(ctx context.Context, path string) (*types.ISBNBook, error) {
	req, err := http.
		NewRequest(k.client).
		Get().
		URL("%s%s", baseURL, path).
		CreateRequest(ctx)
	if err != nil {
		return nil, err
	}

	body, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}

	node, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	return &types.ISBNBook{
		ISBN13: k.getISBNFromBody(node),
		Title:  k.getTitleFromBody(node),
		Volume: k.getVolumeFromPath(path),
		Source: sourceName,
	}, nil
}

func (k kodansha) getVolumeFromPath(path string) string {
	urlSplit := strings.Split(path, "/")
	if len(urlSplit) < 3 {
		return ""
	}
	stubSplit := strings.Split(urlSplit[2], "-")
	volumeMaybe := stubSplit[len(stubSplit)-1]

	if _, err := strconv.Atoi(volumeMaybe); err == nil { // check if it's a number
		return volumeMaybe
	}

	return ""
}

func (k kodansha) getTitleFromBody(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "h2" && utils.AttrEquals(node.Attr, "class", "product-title") {
		if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
			return node.FirstChild.Data
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if title := k.getTitleFromBody(child); len(title) > 0 {
			return title
		}
	}

	return ""
}

func (k kodansha) getISBNFromBody(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "div" && utils.AttrEquals(node.Attr, "class", "product-desktop-rating-table-title-value-wrapper") {
		if isISBN := k.isISBNNode(node); isISBN {
			return k.getISBNFromNode(node)
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isbn := k.getISBNFromBody(child); len(isbn) > 0 {
			return isbn
		}
	}

	return ""
}

func (k kodansha) isISBNNode(node *html.Node) bool {
	if node.Data == "span" && node.FirstChild != nil {
		if strings.HasPrefix(node.FirstChild.Data, "ISBN") {
			return true
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if is := k.isISBNNode(child); is {
			return is
		}
	}

	return false
}

func (k kodansha) getISBNFromNode(node *html.Node) string {
	if node.Data == "span" && node.FirstChild != nil {
		data := node.FirstChild.Data
		if len(data) == 13 {
			if _, err := strconv.Atoi(data); err == nil { // if it's just a number
				return data
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isbn := k.getISBNFromNode(child); len(isbn) > 0 {
			return isbn
		}
	}

	return ""
}
