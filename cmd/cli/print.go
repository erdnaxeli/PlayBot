package main

import (
	"fmt"
	"strings"

	"github.com/erdnaxeli/PlayBot/textbot"
)

func printMusicRecord(result textbot.Result) string {
	var b strings.Builder

	if result.IsNew {
		fmt.Fprint(&b, "+", result.ID, "\n")
	} else {
		fmt.Fprintln(&b, result.ID)
	}
	fmt.Fprintln(&b, result.RecordId)
	fmt.Fprintln(&b, result.Url)
	fmt.Fprintln(&b, "🎉", result.Name)
	fmt.Fprintln(&b, result.Band.Name)
	fmt.Fprintln(&b, result.Duration.Seconds())

	var tags []string
	for _, tag := range result.Tags {
		tags = append(tags, "#"+tag)
	}
	fmt.Fprint(&b, strings.Join(tags, " "))

	if result.Count == 1 {
		fmt.Fprintln(&b, " [1 résultat]")
	} else if result.Count > 1 {
		fmt.Fprintf(&b, " [%d résultats]\n", result.Count)
	} else {
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b, result.Count)

	return b.String()
}
