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
	"errors"
	"log"
	"regexp"
	"sync"

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
	playbot           *playbot.Playbot
	lastCommands      map[types.Channel][]string
	lastCommandsMutex sync.RWMutex
}

func New(playbot *playbot.Playbot) *textBot {
	return &textBot{
		playbot:      playbot,
		lastCommands: make(map[types.Channel][]string),
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
func (t *textBot) Execute(
	channelName string, personName string, msg string,
) (Result, bool, error) {
	channel := types.Channel{Name: channelName}
	person := types.Person{Name: personName}

	args := parseArgs(msg)
	cmd, cmdArgs := args[0], args[1:]
	if cmd == "!" {
		lastCmd, ok := t.getLastCommand(channel)
		if ok {
			log.Printf("Repeat last command: %s", lastCmd)
			cmd, cmdArgs = lastCmd[0], lastCmd[1:]
			args = lastCmd
		}
	}

	var result Result
	ok := true
	var err error

	notImplementedError := errors.New("not implemented")

	switch cmd {
	case "!broken":
		err = notImplementedError
	case "!conf":
		err = notImplementedError
	case "!fav":
		err = t.favCmd(channel, person, cmdArgs)
	case "!later":
		err = notImplementedError
	case "!get":
		result, err = t.getCmd(channel, person, cmdArgs)
	case "!help":
		err = notImplementedError
	case "!stats":
		err = notImplementedError
	case "!tag":
		err = t.saveTagsCmd(channel, person, cmdArgs)
	default:
		result, err = t.saveMusicPost(channel, person, msg)
		ok = false
	}

	if ok {
		t.saveLastCommand(channel, args)
	}

	return result, ok, err
}

func (t *textBot) getLastCommand(channel types.Channel) ([]string, bool) {
	t.lastCommandsMutex.RLock()
	defer t.lastCommandsMutex.RUnlock()

	v, ok := t.lastCommands[channel]
	return v, ok
}

func (t *textBot) saveLastCommand(channel types.Channel, cmd []string) {
	t.lastCommandsMutex.Lock()
	defer t.lastCommandsMutex.Unlock()

	t.lastCommands[channel] = cmd
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
