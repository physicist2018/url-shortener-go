package repoerrors

import "errors"

var (
	ErrURLAlreadyInDB = errors.New("урл уже используется")
)
