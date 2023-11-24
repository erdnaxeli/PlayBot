package matrix

func (c Client) StartSync() error {
	return c.client.Sync()
}
