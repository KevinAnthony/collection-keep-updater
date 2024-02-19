package wikipedia

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/config"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/atye/wikitable2json/pkg/client"
)

type downloader struct {
	client http.Client
}

func NewDownloader(client http.Client) types.CollectionSource {
	if client == nil {
		panic("http client is nil")
	}

	return downloader{
		client: client,
	}
}

func (l downloader) GetISBNs(ctx context.Context, series config.Series) ([]types.ISBNBook, error) {
	tg := client.NewTableGetter("keep-updater")

	tables, err := tg.GetTablesKeyValue(ctx, series.ID, "en", false, 1, series.TableSettings.Table...)
	if err != nil {
		return nil, err
	}

	books := make([]types.ISBNBook, 0, len(tables))
	for _, table := range tables {
		for _, row := range table {
			book := l.processRow(series, row)
			if book != nil {
				books = append(books, *book)
			}
		}
	}
	return books, nil
}

func (l downloader) processRow(series config.Series, row map[string]string) *types.ISBNBook {
	book := types.ISBNBook{
		Volume: l.getVolume(row, series),
		Title:  l.getTitle(row, series),
		ISBN10: l.getISBN10(row, series),
		ISBN13: l.getISBN13(row, series),
	}

	if len(book.Title) == 0 {
		book.Title = fmt.Sprintf("%s Vol %s", series.Name, book.Volume)
	}

	if len(book.ISBN10) > 0 || len(book.ISBN13) > 0 {
		return &book
	}

	return nil
}

func (l downloader) getVolume(row map[string]string, series config.Series) string {
	if series.TableSettings.Volume == nil {
		return ""
	}
	volume, ok := row[*series.TableSettings.Volume]
	if !ok {
		return ""
	}

	v, err := strconv.Atoi(strings.TrimSpace(volume))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%03d", v)
}

func (l downloader) getTitle(row map[string]string, series config.Series) string {
	if series.TableSettings.Title == nil {
		return ""
	}
	title, ok := row[*series.TableSettings.Title]
	if !ok {
		return ""
	}

	return title
}

func (l downloader) getISBN10(row map[string]string, series config.Series) string {
	if series.TableSettings.ISBN == nil {
		return ""
	}

	isbnStr, ok := row[*series.TableSettings.ISBN]
	if !ok {
		return ""
	}

	return l.regexISBN(isbnStr, types.ISBN10regex, 10)
}

func (l downloader) getISBN13(row map[string]string, series config.Series) string {
	if series.TableSettings.ISBN == nil {
		return ""
	}

	isbnStr, ok := row[*series.TableSettings.ISBN]
	if !ok {
		return ""
	}

	return l.regexISBN(isbnStr, types.ISBN13regex, 13)
}

func (l downloader) regexISBN(str string, re *regexp.Regexp, count int) string {
	if re == nil {
		return strings.ReplaceAll(str, "-", "")
	}

	for _, match := range re.FindAllString(str, -1) {
		isbn := strings.ReplaceAll(match, "-", "")
		if len(isbn) == count {
			return isbn
		}
	}

	return ""
}
