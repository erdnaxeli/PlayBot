// Package types implement common types for the project.
package types

// Person represent a user posting posts.
type Person struct {
	Name string
}

// Channel is where a post is done.
type Channel struct {
	Name string
}

// MusicPost is a MusicRecord sent by a Person on a Channel.
type MusicPost struct {
	MusicRecord MusicRecord
	Person
	Channel
}
