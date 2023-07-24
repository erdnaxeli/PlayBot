package types

import "time"

type Band struct {
	Name string
}

type MusicRecord struct {
	Band     Band
	Duration time.Duration
	Name     string
	RecordId string
	Source   string
	Url      string
}
