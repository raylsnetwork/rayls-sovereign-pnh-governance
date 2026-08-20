package withstack

import (
	"fmt"

	cErr "github.com/cockroachdb/errors"
)

type WithStackError struct {
	err error // -> actual WithStack error
}

func (w *WithStackError) Error() string {
	return fmt.Sprintf("%+v", w.err)
}

func (w *WithStackError) Unwrap() error {
	return w.err
}

func Wrap(err error) *WithStackError {
	if err == nil {
		return nil
	}
	return &WithStackError{
		err: cErr.WithStackDepth(err, 1),
	}
}
