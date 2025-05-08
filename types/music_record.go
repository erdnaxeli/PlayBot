package types

import "time"

// Band is the name of the author of a music record.
//
// It can be a single person or it can be a band with multiple people.
type Band struct {
	Name string
}

// MusicRecord contains all the data about a music content.
type MusicRecord struct {
	Band     Band
	Duration time.Duration
	Name     string
	RecordID string
	Source   string
	URL      string
}
