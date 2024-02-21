package types

type VizSettings struct {
	MaximumBacklog *int         `json:"maximum_backlog"`
	Delay          *EncDuration `json:"delay_between"`
}
