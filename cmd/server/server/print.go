package server

import (
	"fmt"
	"strings"

	"github.com/erdnaxeli/PlayBot/textbot"
)

const (
	NORMAL     = "\x0f"
	GREEN      = "\x0303"
	ORANGE     = "\x0307"
	YELLOW     = "\x0308"
	LIGHT_BLUE = "\x0312"
	GREY       = "\x0314"
)

func printMusicRecord(result textbot.Result) string {
	var b strings.Builder

	fmt.Fprint(&b, "🎉 ", YELLOW)

	if result.IsNew {
		fmt.Fprint(&b, "[+", result.ID, "]")
	} else {
		fmt.Fprint(&b, "[", result.ID, "]")
	}

	fmt.Fprint(&b, " ", GREEN, result.Name)

	if result.Band.Name != "" {
		fmt.Fprint(&b, " | ", result.Band.Name)
	}

	fmt.Fprint(&b, " ", LIGHT_BLUE, result.Duration.String())
	if result.Url != "" {
		fmt.Fprint(&b, NORMAL, " => ", result.Url)
	}

	var tags []string
	for _, tag := range result.Tags {
		tags = append(tags, "#"+tag)
	}

	if len(tags) > 0 {
		fmt.Fprint(&b, ORANGE, " ", strings.Join(tags, " "))
	}

	if result.Count == 1 {
		fmt.Fprint(&b, " ", GREY, "[1 résultat]")
	} else if result.Count > 1 {
		fmt.Fprintf(&b, " %s [%d résultats]\n", GREY, result.Count)
	}

	return b.String()
}
