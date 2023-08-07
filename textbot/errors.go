package textbot

type NoRecordFound struct{}

func (e NoRecordFound) Error() string {
	return "no matching record found"
}
