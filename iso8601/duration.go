package iso8601

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

// This regex only parses the time part of an ISO86010  duration expression.
// It is actually invalid as the T (time) delimiter is mandatory, but some
// websites (bandcamp) do not put it.""
var re = regexp.MustCompile(
	`PT?(?:(?P<hours>\d\d?)H)?(?:(?P<minutes>\d\d?)M)(?:(?P<secondes>\d\d?)S)`,
)

// Parse ISO8601 duration.
func ParseDuration(duration string) time.Duration {
	groups := re.FindStringSubmatch(duration)
	if groups == nil {
		log.Fatalf("Invalid duration: %s", duration)
	}

	var hours, minutes, seconds int
	if groups[1] != "" {
		hours, _ = strconv.Atoi(groups[1])
	} else {
		hours = 0
	}

	if groups[2] != "" {
		minutes, _ = strconv.Atoi(groups[2])
	} else {
		minutes = 0
	}

	if groups[3] != "" {
		seconds, _ = strconv.Atoi(groups[3])
	} else {
		seconds = 0
	}

	return time.Duration(hours*int(time.Hour) + minutes*int(time.Minute) + seconds*int(time.Second))
}
