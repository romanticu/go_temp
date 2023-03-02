package funcs

import "errors"

func MakeError() error {
	err := errors.New("i have a err")

	return err
}
