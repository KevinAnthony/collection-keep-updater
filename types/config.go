package types

type Config struct {
	Series    []Series          `json:"series"    mapstructure:"series"`
	Libraries []LibrarySettings `json:"libraries" mapstructure:"libraries"`
}
