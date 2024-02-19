package config

type Series struct {
	Name          string            `json:"name"`
	ID            string            `json:"wikipedia_page"`
	ISBNBlacklist []string          `json:"isbn_blacklist"`
	TableSettings WikiTableSettings `json:"table_settings"`
}

type WikiTableSettings struct {
	Volume *string `json:"volume"`
	Title  *string `json:"title"`
	ISBN   *string `json:"isbn"`
	Table  []int   `json:"tables"`
}
