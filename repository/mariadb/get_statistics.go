package mariadb

import (
	"database/sql"
	"errors"
	"time"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
)

// GetMusicRecordStatistics returns the MusicRecordStatistics corresponding to a given music record.
//
// If the ID references an non existant music record, the zero value of the MusicRecordStatistics type is returned and error is nil.
// You can check this case by looking at MusicRecordStatistics.PostsCount.
// If it is zero, it means there is no post for this music record (and thus the music record does not exist).
func (r Repository) GetMusicRecordStatistics(recordID int64) (playbot.MusicRecordStatistics, error) {
	row := r.db.QueryRow(
		`
		with first_post as (
			select
				pc.sender_irc,
				pc.chan,
				pc.date
			from playbot_chan pc
			where pc.content = ?
			order by pc.date
			limit 1
		),
		posts_count as(
			select
				count(*) as count,
				count(distinct pc.sender_irc) as senders_count,
				count(distinct pc.chan) as chans_count
			from playbot_chan pc
			where pc.content = ?
			group by pc.content
		),
		max_sender as (
			select
				pc.sender_irc,
				count(*) as count
			from playbot_chan pc
			where pc.content = ?
			group by pc.sender_irc
			order by 2 desc
			limit 1
		),
		max_channel as (
			select
				pc.chan,
				count(*) as count
			from playbot_chan pc
			where pc.content = ?
			group by pc.chan
			order by 2 desc
			limit 1
		),
		favorites as (
			select count(*) as count
			from playbot_fav
			where id = ?
		)
		select
			first_post.sender_irc,
			first_post.chan,
			first_post.date,
			posts_count.count,
			posts_count.senders_count,
			posts_count.chans_count,
			max_sender.sender_irc,
			max_sender.count,
			max_channel.chan,
			max_channel.count,
			favorites.count
		from
			first_post,
			posts_count,
			max_sender,
			max_channel,
			favorites
		`,
		recordID,
		recordID,
		recordID,
		recordID,
		recordID,
	)

	var firstPersonName string
	var firstChannelName string
	var firstPostDate time.Time
	var postsCount int64
	var sendersCount int64
	var channelsCount int64
	var maxPersonName string
	var maxPersonCount int64
	var maxChannelName string
	var maxChannelCount int64
	var favoritesCount int64
	err := row.Scan(
		&firstPersonName,
		&firstChannelName,
		&firstPostDate,
		&postsCount,
		&sendersCount,
		&channelsCount,
		&maxPersonName,
		&maxPersonCount,
		&maxChannelName,
		&maxChannelCount,
		&favoritesCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return playbot.MusicRecordStatistics{}, nil
		}

		return playbot.MusicRecordStatistics{}, err
	}

	return playbot.MusicRecordStatistics{
		PostsCount:       int(postsCount),
		PeopleCount:      int(sendersCount),
		ChannelsCount:    int(channelsCount),
		MaxPerson:        types.Person{Name: maxPersonName},
		MaxPersonCount:   int(maxPersonCount),
		MaxChannel:       types.Channel{Name: maxChannelName},
		MaxChannelCount:  int(maxChannelCount),
		FirstPostPerson:  types.Person{Name: firstPersonName},
		FirstPostChannel: types.Channel{Name: firstChannelName},
		FirstPostDate:    firstPostDate,
		FavoritesCount:   int(favoritesCount),
	}, nil
}
