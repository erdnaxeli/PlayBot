package textbot

import (
	"context"
	"log"
	"strings"

	"github.com/erdnaxeli/PlayBot/types"
	"github.com/spf13/pflag"
)

func (t textBot) getCmd(channel types.Channel, person types.Person, args []string) (Result, error) {
	var all bool
	var force bool

	flagSet := pflag.NewFlagSet("!get", pflag.ContinueOnError)
	flagSet.BoolVarP(
		&force,
		"force",
		"f",
		false,
		"retourne des résultats même s'il y en a beaucoup",
	)
	flagSet.BoolVarP(
		&all,
		"all",
		"a",
		false,
		"recherche dans tous les contenus, pas seulement ceux du channel courant",
	)

	err := flagSet.Parse(args)
	if err != nil {
		return Result{}, err
	}

	var searchChannel types.Channel
	if !all {
		searchChannel = channel
	}

	var words []string
	var tags []string
	for _, arg := range flagSet.Args() {
		if strings.HasPrefix(arg, "#") {
			tags = append(tags, arg[1:])
		} else {
			words = append(words, arg)
		}
	}

	ctx := context.Background()
	count, ch, err := t.playbot.SearchMusicRecord(ctx, searchChannel, words, tags)
	if err != nil {
		return Result{}, err
	}

	searchResult, ok := <-ch
	log.Print(searchResult)
	if !ok {
		return Result{}, NoRecordFound{}
	}

	resultTags, err := t.playbot.GetTags(searchResult.Id())
	if err != nil {
		return Result{}, nil
	}

	result := Result{
		ID:          searchResult.Id(),
		MusicRecord: searchResult.MusicRecord(),
		Tags:        resultTags,
		IsNew:       false,
		Count:       count,
	}
	return result, nil
}
