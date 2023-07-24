package types

type Person struct {
	Name string
}

type Channel struct {
	Name string
}

type MusicPost struct {
	MusicRecord
	Person
	Channel
}
