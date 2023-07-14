package iso8601

import (
	"regexp"
	"strconv"
	"time"
)

// Parse ISO8601 duration.
func ParseDuration(duration string) time.Duration {
	re := regexp.MustCompile(`PT(?:(?P<hours>\d\d?)H)?(?:(?P<minutes>\d\d?)M)(?:(?P<secondes>\d\d?)S)`)
	groups := re.FindStringSubmatch(duration)

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
