package textbot

import (
	"fmt"

	"github.com/erdnaxeli/PlayBot/types"
)

// !fav is the only command that can be used with an url. It saves the post and adds it
// to the user favorites in one command.
// We need to catch the !fav command to stop the process and not save the eventual post
// without adding it into the user favorites, which will be misleading.
func (t textBot) favCmd(channel types.Channel, person types.Person, args []string) error {
	return fmt.Errorf("not implemented")
}
