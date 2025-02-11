package repoerrors

import "errors"

var (
	ErrUrlAlreadyInDB = errors.New("урл уже используется")
)
