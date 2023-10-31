package model

type ErrNotFound struct {
	Resource string
}

func (e *ErrNotFound) Error() string {
	return e.Resource + "Not Found"
}
