package wikipedia

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/atye/wikitable2json/pkg/client"
)

type wikiSource struct {
	client http.Client
}

func New(client http.Client) types.ISource {
	if client == nil {
		panic("http client is nil")
	}

	return wikiSource{
		client: client,
	}
}

func (l wikiSource) GetISBNs(ctx context.Context, series types.Series) (types.ISBNBooks, error) {
	if true {
		return nil, nil
	}

	tg := client.NewTableGetter("keep-updater")
	settings, ok := series.SourceSettings.(types.WikipediaSettings)
	if !ok {
		return nil, fmt.Errorf("setting type not correct")
	}

	tables, err := tg.GetTablesKeyValue(ctx, series.ID, "en", false, 1, settings.Table...)
	if err != nil {
		return nil, err
	}

	books := types.NewISBNBooks(len(tables))
	for _, table := range tables {
		for _, row := range table {
			book := l.processRow(series, settings, row)
			if book != nil {
				books = append(books, *book)
			}
		}
	}
	return books, nil
}

func (l wikiSource) processRow(series types.Series, settings types.WikipediaSettings, row map[string]string) *types.ISBNBook {
	book := types.ISBNBook{
		Volume: l.getVolume(row, settings),
		Title:  l.getTitle(row, settings),
		ISBN10: l.getISBN10(row, settings),
		ISBN13: l.getISBN13(row, settings),
	}

	if len(book.Title) == 0 {
		book.Title = fmt.Sprintf("%s Vol %s", series.Name, book.Volume)
	}

	if len(book.ISBN10) > 0 || len(book.ISBN13) > 0 {
		return &book
	}

	return nil
}

func (l wikiSource) getVolume(row map[string]string, tableSetting types.WikipediaSettings) string {
	if tableSetting.Volume == nil {
		return ""
	}
	volume, ok := row[*tableSetting.Volume]
	if !ok {
		return ""
	}

	v, err := strconv.Atoi(strings.TrimSpace(volume))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%03d", v)
}

func (l wikiSource) getTitle(row map[string]string, tableSetting types.WikipediaSettings) string {
	if tableSetting.Title == nil {
		return ""
	}
	title, ok := row[*tableSetting.Title]
	if !ok {
		return ""
	}

	return title
}

func (l wikiSource) getISBN10(row map[string]string, tableSetting types.WikipediaSettings) string {
	if tableSetting.ISBNColumnTitle == nil {
		return ""
	}

	isbnStr, ok := row[*tableSetting.ISBNColumnTitle]
	if !ok {
		return ""
	}

	return l.regexISBN(isbnStr, types.ISBN10regex, 10)
}

func (l wikiSource) getISBN13(row map[string]string, tableSetting types.WikipediaSettings) string {
	if tableSetting.ISBNColumnTitle == nil {
		return ""
	}

	isbnStr, ok := row[*tableSetting.ISBNColumnTitle]
	if !ok {
		return ""
	}

	return l.regexISBN(isbnStr, types.ISBN13regex, 13)
}

func (l wikiSource) regexISBN(str string, re *regexp.Regexp, count int) string {
	if re == nil {
		return utils.ISBNNormalize(str)
	}

	for _, match := range re.FindAllString(str, -1) {
		isbn := utils.ISBNNormalize(match)
		if len(isbn) == count {
			return isbn
		}
	}

	return ""
}
