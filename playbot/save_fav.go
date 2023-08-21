package playbot

func (p *Playbot) SaveFav(user string, recordID int64) error {
	return p.repository.SaveFav(user, recordID)
}
