package play

import "fmt"

func derr(err *error, msg string) {
	if err != nil && *err != nil {
		*err = fmt.Errorf("%s: %w", msg, *err)
	}
}
