package types

import (
	"encoding/json"
	"time"

	"github.com/kevinanthony/collection-keep-updater/utils"
)

type EncDuration struct {
	time.Duration
}

func (d *EncDuration) UnmarshalText(text []byte) error {
	var err error
	parsed, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	d.Duration = parsed

	return nil
}

func (d *EncDuration) UnmarshalJSON(data []byte) error {
	rm := json.RawMessage(data)

	text, err := utils.Unmarshal[string](&rm)
	if err != nil {
		return err
	}

	parsed, err := time.ParseDuration(text)
	if err != nil {
		return err
	}

	d.Duration = parsed

	return nil
}
