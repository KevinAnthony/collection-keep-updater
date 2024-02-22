package types

import "github.com/kevinanthony/collection-keep-updater/utils"

type WikipediaSettings struct {
	Volume          *string `json:"volume"`
	Title           *string `json:"title"`
	ISBNColumnTitle *string `json:"isbn_column_title"`
	Table           []int   `json:"tables"`
}

func newWikipediaSettings(data map[string]interface{}) *WikipediaSettings {
	if len(data) == 0 {
		return nil
	}

	settings := WikipediaSettings{
		Volume:          utils.GetPtr[string](data, "volume"),
		Title:           utils.GetPtr[string](data, "title"),
		ISBNColumnTitle: utils.GetPtr[string](data, "isbn_column_title"),
		Table:           utils.GetArray[int](data, "tables"),
	}

	if settings.ISBNColumnTitle == nil || len(*settings.ISBNColumnTitle) == 0 {
		return nil
	}

	return &settings
}
