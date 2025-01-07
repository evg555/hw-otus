package storage

import "errors"

var (
	ErrStorageNotExist    = errors.New("storage not exist")
	ErrDateBusy           = errors.New("date is busy by another event")
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventNotExists     = errors.New("event not exist")
)
