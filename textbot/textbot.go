// Implementation of a Playbot using text message to interact.
//
// Commands are used with an exclamation mark. Currently implemented commands are:
// * !get
//
// If no supported command is found, it looks for an URL in the message and try to save
// the corresponding music record. Tags can be added alongside the URL with "#":
// > this #mix is so awesome https://soundcloud.com/hate_music/frederic-hate-podcast-332 #techno
package textbot

import (
	"regexp"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
)

// Result represent the result of a command or the saved music record.
type Result struct {
	ID int64
	types.MusicRecord
	// The tags of the MusicRecord.
	Tags []string
	// True if the MusicRecord was newly inserted, false else.
	IsNew bool
	// If the MusicRecord come from a search, this is the count of results.
	Count int64
}

type textBot struct {
	playbot playbot.Playbot
}

func New(playbot playbot.Playbot) textBot {
	return textBot{
		playbot: playbot,
	}
}

// Execute try to parse the given message and execute the found command or save the
// found music record.
//
// If a command is found, the bool returned value is true and the Result returned value
// contains the music record returned by the command, if any.
// If no command is found, the bool returned value is false and the Result returned
// value contains the music record saved, if any.
// If the Result object is equal to its zero value and the bool value is false, it means
// nothing has been done.
func (t textBot) Execute(
	channelName string, personName string, msg string,
) (Result, bool, error) {
	channel := types.Channel{Name: channelName}
	person := types.Person{Name: personName}

	args := parseArgs(msg)
	cmd, args := args[0], args[1:]
	switch cmd {
	case "!fav":
		err := t.favCmd(channel, person, args)
		return Result{}, true, err
	case "!get":
		result, err := t.getCmd(channel, person, args)
		return result, true, err
	default:
		result, err := t.saveMusicPost(channel, person, msg)
		return result, false, err
	}

}

func parseArgs(msg string) []string {
	blankRgx := regexp.MustCompile(`\s+`)
	args := blankRgx.Split(msg, -1)

	cleanArgs := args[:0]
	for _, v := range args {
		if v != "" {
			cleanArgs = append(cleanArgs, v)
		}
	}

	return cleanArgs
}
