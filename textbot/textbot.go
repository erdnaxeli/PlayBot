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
	"log"
	"regexp"
	"strconv"
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
	// Some command may return statistics
	Statistics playbot.MusicRecordStatistics
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
	channelName string, personName string, msg string, user string,
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

	switch cmd {
	case "!broken":
		err = NotImplementedError
	case "!conf":
		err = NotImplementedError
	case "!fav":
		result, err = t.favCmd(channel, person, cmdArgs, user)
	case "!later":
		err = NotImplementedError
	case "!get":
		result, err = t.getCmd(channel, person, cmdArgs)
	case "!help":
		err = NotImplementedError
	case "!stats":
		result, err = t.statsCmd(channel, person, cmdArgs)
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

func (t *textBot) getRecordIDFromArgs(channel types.Channel, args []string) (int64, []string, error) {
	recordID, args := parseID(args)
	if recordID > 0 {
		return recordID, args, nil
	}

	if recordID < -10 {
		return 0, args, OffsetToBigError
	}

	recordID, err := t.playbot.GetLastID(channel, int(recordID))
	if err != nil {
		return 0, args, err
	}

	return recordID, args, nil
}

func parseID(args []string) (int64, []string) {
	if len(args) == 0 {
		return 0, args
	}

	recordID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return 0, args
	}

	return recordID, args[1:]
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
