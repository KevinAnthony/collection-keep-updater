package types

import (
	"encoding/json"
	"fmt"

	"github.com/kevinanthony/collection-keep-updater/utils"

	"github.com/pkg/errors"
)

type Series struct {
	Name           string     `json:"name"`
	ID             string     `json:"id"`
	ISBNBlacklist  []string   `json:"isbn_blacklist"`
	Source         SourceType `json:"source"`
	SourceSettings any        `json:"source_settings"`
}

func (s *Series) UnmarshalJSON(data []byte) error {
	var raw map[string]*json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	name, err := utils.Unmarshal[string](raw["name"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series name")
	}

	id, err := utils.Unmarshal[string](raw["id"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series id")
	}

	source, err := utils.Unmarshal[SourceType](raw["source"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series source")
	}

	blacklist, err := utils.Unmarshal[[]string](raw["isbn_blacklist"])
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal series blacklist")
	}

	s.Name = name
	s.ID = id
	s.Source = source
	s.ISBNBlacklist = blacklist

	switch s.Source {
	case WikipediaSource:
		settings, err := utils.Unmarshal[WikipediaSettings](raw["source_settings"])
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal source_settings to wikipedia")
		}

		s.SourceSettings = settings
	case VizSource:
		settings, err := utils.Unmarshal[VizSettings](raw["source_settings"])
		if err != nil {
			return errors.Wrap(err, "unable to unmarshal source_settings to viz")
		}

		s.SourceSettings = settings
	default:
		return fmt.Errorf("unknown source type: %s", s.Source)
	}

	return nil
}
