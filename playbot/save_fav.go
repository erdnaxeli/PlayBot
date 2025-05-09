package playbot

// SaveFav save a music record to a user's favorites.
func (p *Playbot) SaveFav(user string, recordID int64) error {
	return p.repository.SaveFav(user, recordID)
}
