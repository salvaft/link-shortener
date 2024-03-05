package persistance

type URLNotFound struct {
	Message string
}

func (ce *URLNotFound) Error() string {
	return ce.Message
}
