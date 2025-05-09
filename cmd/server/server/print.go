package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
)

// Constant holding strinrgs to change the color of a IRC message.
const (
	Normal    = "\x0f"
	Green     = "\x0303"
	Orange    = "\x0307"
	Yellow    = "\x0308"
	LightBlue = "\x0312"
	Grey      = "\x0314"
)

// IrcMusicRecordPrinter provides the PrintResult method.
type IrcMusicRecordPrinter struct{}

// PrintResult formats a string according to the given Result.
func (IrcMusicRecordPrinter) PrintResult(result textbot.Result) string {
	var b strings.Builder

	fmt.Fprint(&b, Yellow)

	if result.IsNew {
		fmt.Fprint(&b, "[+", result.ID, "]")
	} else {
		fmt.Fprint(&b, "[", result.ID, "]")
	}

	fmt.Fprint(&b, " ", Green, result.Name)

	if result.Band.Name != "" {
		fmt.Fprint(&b, " | ", result.Band.Name)
	}

	fmt.Fprint(&b, " ", LightBlue, result.Duration.String())
	if result.URL != "" {
		fmt.Fprint(&b, Normal, " => ", result.URL)
	}

	var tags []string
	for _, tag := range result.Tags {
		tags = append(tags, "#"+tag)
	}

	if len(tags) > 0 {
		fmt.Fprint(&b, Orange, " ", strings.Join(tags, " "))
	}

	if result.Count == 1 {
		fmt.Fprint(&b, " ", Grey, "[1 résultat]")
	} else if result.Count > 1 {
		fmt.Fprintf(&b, " %s [%d résultats]\n", Grey, result.Count)
	}

	return b.String()
}

// IrcStatisticsPrinter provides the PrintStatistics method.
type IrcStatisticsPrinter struct {
	// The location is which times will be shown.
	Location *time.Location
}

// PrintStatistics formats the a string according to the statistics object given.
func (s IrcStatisticsPrinter) PrintStatistics(statistics playbot.MusicRecordStatistics) string {
	var resultMsg strings.Builder
	fmt.Fprintf(
		&resultMsg,
		"Posté la 1re fois par %s le %s sur %s.",
		statistics.FirstPostPerson.Name,
		statistics.FirstPostDate.In(s.Location).Format("02-01-2006 15:04:05"),
		statistics.FirstPostChannel.Name,
	)

	fmt.Fprintf(
		&resultMsg,
		" Posté %d fois par %d personne",
		statistics.PostsCount,
		statistics.PeopleCount,
	)
	plural(statistics.PeopleCount, &resultMsg)

	fmt.Fprintf(
		&resultMsg, " sur %d channel", statistics.ChannelsCount,
	)
	plural(statistics.ChannelsCount, &resultMsg)
	fmt.Fprint(&resultMsg, ".")

	fmt.Fprintf(
		&resultMsg,
		" %s l'a posté %d fois.",
		statistics.MaxPerson.Name,
		statistics.MaxPersonCount,
	)

	fmt.Fprintf(
		&resultMsg,
		" Sauvegardé dans ses favoris par %d personne",
		statistics.FavoritesCount,
	)
	plural(statistics.FavoritesCount, &resultMsg)
	fmt.Fprint(&resultMsg, ".")

	return resultMsg.String()
}
