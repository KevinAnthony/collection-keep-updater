package types

const WikipediaSource SourceType = "wikipedia"

type WikipediaSettings struct {
	Volume *string `json:"volume"`
	Title  *string `json:"title"`
	ISBN   *string `json:"isbn"`
	Table  []int   `json:"tables"`
}
