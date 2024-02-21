package types

type WikipediaSettings struct {
	Volume          *string `json:"volume"`
	Title           *string `json:"title"`
	ISBNColumnTitle *string `json:"isbn_column_title"`
	Table           []int   `json:"tables"`
}
