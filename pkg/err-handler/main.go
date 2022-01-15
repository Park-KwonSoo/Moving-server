package errhandler

import (
	"errors"
)

func NotFoundErr() error {
	return errors.New("Not Found Error")
}
