package common

import (
	"errors"
	"fmt"
)

// ArgParse takes in arguments passed by
// the caller, *excluding the first argument*,
// and returns a Conf reflecting their choices.
// If arguments are invalid, returns an error.
// This could serve to be polished up a bit, but
// it works for now and is somewhat flexible.
func ArgParse(a []string) (Conf, error) {
	i := 0
	c := MkConfig()
	chose := false

	for i < len(a) {
		if chose {
			return c, errors.New("usage")
		}

		if a[i] == "-s" {
			c.Serving = true
			c.Name = c.MkName()
			chose = true
			i++
		} else if a[i] == "-c" {
			c.Serving = false
			i++

			if i == len(a) {
				return c, errors.New("usage")
			}

			nm, err := c.IDToName(a[i])
			if err != nil {
				fmt.Println("Session name not found!")
				return c, errors.New("file not found")
			}

			c.Name = nm
			chose = true
			i++

		} else {
			return c, errors.New("usage")
		}
	}

	if chose {
		return c, nil
	}
	return c, errors.New("usage")
}
