# PlayBot

This is a rewrite of [this project](https://git.iiens.net/morignot2011/playbot.old), a Perl IRC bot.

## The bot

The bot is used to share music.
You can post a link to a supported website and it will:
* extract the music content
* save it in a database
* save any additional tags given with `#`

You can then:
* save a music into your favorites with `!fav`
* search a music with `!get some query #with #tags` (or get a random music if called without search parameters)
* get info about a music: when it was shared for the first time, by whom, how many times it was shared again

See this example of interaction:
![example of interaction with the PlayBot](irc.png)

The second message from the bot is telling that it cannot find anything in the current channel with those search parameters but that results have been found using a global search.

A website also exists to go through a channel history or see your personal favorites.
This part will need a complete rewrite and rethinking from scratch.

## The rewrite

The goal is to do a progressive roll out between the old Perl bot and the new golang bot.
In order to do that, the Perl bot has been modified to execute the golang bot executable with the received message as parameter.
If the process end with success the Perl bot uses the data returned.
Else it uses its own extractors.

With that we can implement the extractors one by one in golang, and the Perl bot will automagically use them.
The extractors currently implemented are:
* SoundCloud and YouTube: the two websites mainly used by the bot users
* Bandcamp: a new extractor not supported by the Perl bot (already a new feature, yay!)

### Next steps

A next step will be to implement the music save into the database.
This will need to edit the Perl bot to not save the data when the golang process succeed.
Note that this means saving the music data **and** any tags added by the user (so wit need to be able to get the tags from the input message).

When this will be implemented it would be easier to iterate.
The following step will be to implement the parsing of commands in the form `!somecommand`, and then to implement the main commands:
* `!tag` to add tag to a previous music
* `!get` to do a search

This will need to edit the Perl bot to execute the golang process earlier, before parsing any commands from the user input.

The `!fav` command will be more difficult to implement as it currently needs a process with private messages to verify the user authentication (against NickServ).
And this feature implementation would probably need a rework if we drop IRC.

Because yes, the idea is to never implement IRC communication in the golang bot.
Instead the first protocol to be implemented will be [[matrix]](https://matrix.org).
As matrix supports bridges to other communication protocols, and especially IRC, the current users will not be dropped.
Implementing other protocols like Slack or Discord could be an idea in the future if the demand is here.

And finally, once the bot is able to talk directly to the users, the Perl bot will not be needed anymore.
We could then host the bot (and the database) elsewhere, then migrate the database from MariaDB to PostgreSQL.
This will also implies to move the website somewhere else to, before rewriting it.

## Hosting an instance of the bot

As you can see, the following README is about the current rewrite and the current running instance, which you will probably never use (the IRC server it runs on is public but small).

That's because my first goal for this project is to not break the current running bot for its users (which includes myself). When the rewrite will be done and the current instance swapped with this new implementation I will edit this README with instruction on how to host it.
