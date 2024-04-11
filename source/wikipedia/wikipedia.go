package wikipedia

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/atye/wikitable2json/pkg/client"
)

//go:generate mockery --srcpkg=github.com/atye/wikitable2json/pkg/client --name=TableGetter --structname=TableGetterMock --filename=table_getter_mock.go --output . --outpkg=wikipedia

const (
	sourceName = "Wikipedia"
)

type wikiSource struct {
	settingsHelper
	client      http.Client
	tableGetter client.TableGetter
}

func New(cmd types.ICommand) types.ISource {
	return wikiSource{
		settingsHelper: settingsHelper{},
		client:         ctxu.GetHttpClient(cmd),
		tableGetter:    ctxu.GetWikiGetter(cmd),
	}
}

func (l wikiSource) GetISBNs(ctx context.Context, series types.Series) (types.ISBNBooks, error) {
	settings, err := types.GetSetting[wikiSettings](series)
	if err != nil {
		return nil, err
	}

	tables, err := l.tableGetter.GetTablesKeyValue(ctx, series.ID, "en", false, 1, settings.Table...)
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

	for _, blackISBN := range series.ISBNBlacklist {
		index := books.FindIndexByISBN(blackISBN)
		if index >= 0 {
			books = books.RemoveAt(index)
		}
	}

	return books, nil
}

func (l wikiSource) processRow(series types.Series, settings wikiSettings, row map[string]string) *types.ISBNBook {
	book := types.ISBNBook{
		Volume: l.getVolume(row, settings),
		Title:  l.getTitle(row, settings),
		ISBN10: l.getISBN10(row, settings),
		ISBN13: l.getISBN13(row, settings),
		Source: sourceName,
	}

	if len(book.Title) == 0 {
		book.Title = fmt.Sprintf("%s Vol %s", series.Name, book.Volume)
	}

	if len(book.ISBN10) > 0 || len(book.ISBN13) > 0 {
		return &book
	}

	return nil
}

func (l wikiSource) getVolume(row map[string]string, tableSetting wikiSettings) string {
	if tableSetting.VolumeHeader == nil {
		return ""
	}
	volume, ok := row[*tableSetting.VolumeHeader]
	if !ok {
		return ""
	}

	v, err := strconv.Atoi(strings.TrimSpace(volume))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%03d", v)
}

func (l wikiSource) getTitle(row map[string]string, tableSetting wikiSettings) string {
	if tableSetting.TitleHeader == nil {
		return ""
	}
	title, ok := row[*tableSetting.TitleHeader]
	if !ok {
		return ""
	}

	return title
}

func (l wikiSource) getISBN10(row map[string]string, tableSetting wikiSettings) string {
	if tableSetting.ISBNHeader == nil {
		return ""
	}

	isbnStr, ok := row[*tableSetting.ISBNHeader]
	if !ok {
		return ""
	}

	return l.regexISBN(isbnStr, types.ISBN10regex, 10)
}

func (l wikiSource) getISBN13(row map[string]string, tableSetting wikiSettings) string {
	if tableSetting.ISBNHeader == nil {
		return ""
	}

	isbnStr, ok := row[*tableSetting.ISBNHeader]
	if !ok {
		return ""
	}

	return l.regexISBN(isbnStr, types.ISBN13regex, 13)
}

func (l wikiSource) regexISBN(str string, re *regexp.Regexp, count int) string {
	for _, match := range re.FindAllString(str, -1) {
		isbn := utils.ISBNNormalize(match)
		if len(isbn) == count {
			return isbn
		}
	}

	return ""
}
