package types

type LibrarySettings struct {
	Name        LibraryType `json:"type"                 mapstructure:"type"`
	WantedColID string      `json:"wanted_collection_id" mapstructure:"wanted_collection_id"`
	OtherColIDs []string    `json:"other_collection_ids" mapstructure:"other_collection_ids"`
	APIKey      string      `json:"api_key"              mapstructure:"api_key"`
}
