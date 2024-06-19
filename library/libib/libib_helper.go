package libib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"slices"
	"strings"

	"github.com/kevinanthony/collection-keep-updater/ctxu"
	"github.com/kevinanthony/collection-keep-updater/types"
	"github.com/kevinanthony/collection-keep-updater/utils"
	"github.com/kevinanthony/gorps/v2/http"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

//go:generate mockery --name=ILibibHelper --structname=ILibibHelperMock --filename=libib_helper_mock.go --inpackage
type ILibibHelper interface {
	SubmitForm(cmd types.ICommand, wanted types.ISBNBooks) (submitResults, error)
	GetISBNFromSuccessResult(cmd types.ICommand, results submitResults) ([]string, error)
	ValidateResults(cmd types.ICommand, wanted types.ISBNBooks, submittedISBNs []string) error
	SaveResults(cmd types.ICommand, s []string) error
	GetQueryParamsFromSuccessResults(cmd types.ICommand, results submitResults) ([]string, error)
}

type libibHelper struct {
	cfg    types.LibrarySettings
	client http.Client
}

func NewLibinHelper(cmd types.ICommand, cfg types.LibrarySettings) ILibibHelper {
	return libibHelper{
		client: ctxu.GetHttpClient(cmd),
		cfg:    cfg,
	}
}

func (l libibHelper) SubmitForm(cmd types.ICommand, wanted types.ISBNBooks) (submitResults, error) {
	csv, err := gocsv.MarshalString(l.createCSVEntries(wanted))
	if err != nil {
		return submitResults{}, err
	}
	payload := &bytes.Buffer{}
	form := multipart.NewWriter(payload)
	_ = form.WriteField("csv-import-library-select", l.cfg.WantedColID)
	_ = form.WriteField("csv-import-type", "book")

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="csv-import-file-0"; filename="wanted.csv"`)
	h.Set("Content-Type", "application/octet-stream")
	writer, _ := form.CreatePart(h)
	_, _ = io.Copy(writer, strings.NewReader(csv))
	_ = form.Close()

	req, _ := http.
		NewRequest(l.client).
		Post().
		URL(validateURL).
		Body(payload.String()).
		Header("Cookie", l.cfg.APIKey).
		Header("Content-Type", form.FormDataContentType()).
		CreateRequest(cmd.Context())
	var result submitResults
	body, err := l.client.Do(req)
	if err != nil {
		// TODO: log error
		return submitResults{}, err
	}

	bdiddy, _ := io.ReadAll(body)

	err = json.Unmarshal(bdiddy, &result)
	if err != nil {
		return submitResults{}, err
	}

	if result.Outcome != "success" {
		return submitResults{}, errors.New("unable to upload CSV.  success = false")
	}

	return result, nil
}

func (l libibHelper) GetISBNFromSuccessResult(_ types.ICommand, results submitResults) ([]string, error) {
	node, err := html.Parse(strings.NewReader(results.CSVPreview))
	if err != nil {
		return nil, err
	}

	return l.walkISBNResults(node), nil
}

func (l libibHelper) GetQueryParamsFromSuccessResults(cmd types.ICommand, results submitResults) ([]string, error) {
	node, err := html.Parse(strings.NewReader(results.CSVPreview))
	if err != nil {
		return nil, err
	}

	return l.walkParameterResults(node), nil
}

func (l libibHelper) ValidateResults(cmd types.ICommand, wanted types.ISBNBooks, submittedISBNs []string) error {
	for _, isbn := range wanted {
		if !slices.Contains(submittedISBNs, isbn.ISBN13) && !slices.Contains(submittedISBNs, isbn.ISBN10) {
			return fmt.Errorf("isbn not found: %s/%s", isbn.ISBN10, isbn.ISBN13)
		}
	}

	return nil
}

func (l libibHelper) createCSVEntries(books types.ISBNBooks) []libibCSVEntries {
	entries := make([]libibCSVEntries, 0, len(books))
	for _, book := range books {
		entry := libibCSVEntries{
			ISBN13: book.ISBN13,
			ISBN:   book.ISBN10,
		}

		entries = append(entries, entry)
	}

	return entries
}

func (l libibHelper) walkISBNResults(node *html.Node) []string {
	var results []string

	if node.Type == html.ElementNode && node.Data == "strong" {
		if node.FirstChild != nil &&
			(node.FirstChild.Data == "ean_isbn13" || node.FirstChild.Data == "upc_isbn10") {
			for child := node; child != nil; child = child.NextSibling {
				if child.Data != "ul" {
					continue
				}

				for list := child.FirstChild; list != nil; list = list.NextSibling {
					if list.FirstChild != nil {
						results = append(results, list.FirstChild.Data)
					}
				}
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if isbn := l.walkISBNResults(child); len(isbn) != 0 {
			return append(results, isbn...)
		}
	}

	return results
}

func (l libibHelper) walkParameterResults(node *html.Node) []string {
	var colTitle []string

	if node.Type == html.ElementNode && node.Data == "strong" {
		if countChildren(node.Parent) == 3 {
			valueSelected := getSelected(node.NextSibling)
			if len(valueSelected) > 0 {
				id, found := utils.AttrContains(node.Parent.Attr, "data-column-id")
				if found {
					colTitle = append(colTitle, fmt.Sprintf("col_obj[%s]=%s", id, valueSelected))
				}
			}
		}
	}
	if node.Data == "input" && utils.AttrEquals(node.Attr, "type", "hidden") {
		id, idFound := utils.AttrContains(node.Attr, "id")
		value, valueFound := utils.AttrContains(node.Attr, "value")
		if idFound && valueFound {
			colTitle = append(colTitle, fmt.Sprintf("%s=%s", id, value))
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if col := l.walkParameterResults(child); len(col) > 0 {
			colTitle = append(colTitle, col...)
		}
	}

	return colTitle
}

func getSelected(node *html.Node) string {
	for child := node.FirstChild.FirstChild; child != nil; child = child.NextSibling {
		if _, found := utils.AttrContains(child.Attr, "selected"); found {
			if value, found := utils.AttrContains(child.Attr, "value"); found && len(value) > 0 {
				return value
			}
		}
	}

	return ""
}

func countChildren(node *html.Node) int {
	i := 0
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		i++
	}

	return i
}

func (l libibHelper) SaveResults(cmd types.ICommand, params []string) error {
	bodyStr := strings.Join(params, "&")

	req, _ := http.
		NewRequest(l.client).
		Post().
		URL(processURL).
		Body(bodyStr).
		Header("Cookie", l.cfg.APIKey).
		Header("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8").
		CreateRequest(cmd.Context())
	body, err := l.client.Do(req)
	if err != nil {
		return err
	}
	bts, _ := io.ReadAll(body)
	fmt.Println(string(bts))
	return nil
}

type submitResults struct {
	Outcome    string `json:"outcome"`
	CSVPreview string `json:"csv-preview"`
}
