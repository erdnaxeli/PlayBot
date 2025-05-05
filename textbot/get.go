package textbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/spf13/pflag"
)

func (t *textBot) getCmd(channel types.Channel, _ types.Person, args []string) (Result, error) {
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

	var recordID int64
	var record types.MusicRecord
	var count int64

	recordID, record, err = t.getByID(flagSet.Args())
	if err != nil {
		return Result{}, err
	} else if recordID == 0 {
		var words []string
		var tags []string
		var excludedTags []string

		for _, arg := range flagSet.Args() {
			if strings.HasPrefix(arg, "##") {
				excludedTags = append(excludedTags, arg[2:])
			} else if strings.HasPrefix(arg, "#") {
				tags = append(tags, arg[1:])
			} else {
				words = append(words, arg)
			}
		}

		ctx := context.Background()
		//nolint:govet
		searchContext, _ := context.WithTimeout(ctx, 6*time.Hour)
		resultContext, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var searchResult playbot.SearchResult
		count, searchResult, err = t.playbot.SearchMusicRecord(
			resultContext,
			playbot.Search{
				Ctx:          searchContext,
				Channel:      channel,
				GlobalSearch: all,
				Words:        words,
				Tags:         tags,
				ExcludedTags: excludedTags,
			},
		)
		if err != nil {
			if errors.Is(err, playbot.SearchCanceledError{}) {
				// we retry once
				log.Print("The search was canceled, we retry once")
				count, searchResult, err = t.playbot.SearchMusicRecord(
					resultContext,
					playbot.Search{
						Ctx:          searchContext,
						Channel:      channel,
						GlobalSearch: all,
						Words:        words,
						Tags:         tags,
					},
				)
				if err != nil {
					if errors.Is(err, playbot.SearchCanceledError{}) {
						return Result{}, fmt.Errorf("the search keeps timeouting: %w", err)
					}

					return Result{Count: count}, err
				}
			} else {
				return Result{Count: count}, err
			}
		}

		recordID = searchResult.ID()
		record = searchResult.MusicRecord()
	}

	resultTags, err := t.playbot.GetTags(recordID)
	if err != nil {
		return Result{}, err
	}

	err = t.playbot.SaveMusicPost(recordID, channel, types.Person{Name: "PlayBot"})
	if err != nil {
		return Result{}, err
	}

	result := Result{
		ID:          recordID,
		MusicRecord: record,
		Tags:        resultTags,
		IsNew:       false,
		Count:       count,
	}
	return result, nil
}

func (t *textBot) getByID(args []string) (int64, types.MusicRecord, error) {
	if len(args) != 1 {
		return 0, types.MusicRecord{}, nil
	}

	recordID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		if _, ok := err.(*strconv.NumError); ok {
			// not a number
			return 0, types.MusicRecord{}, nil
		}

		// unknown error
		return 0, types.MusicRecord{}, err
	}

	record, err := t.playbot.GetMusicRecord(int64(recordID))

	return recordID, record, err
}
