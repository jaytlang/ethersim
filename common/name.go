package common

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// MkName makes a name suitable for a new UDS connection
// via the RNG and a little bit of error checking, along
// with the given configuration structure
func (c *Conf) MkName() string {
	var nm string
	var id int

	rand.Seed(time.Now().UnixNano())

gen:
	id = rand.Intn(c.MaxID-c.MinID) + c.MinID
	nm = fmt.Sprintf("%s/%d", c.Prefix, id)
	if _, err := os.Stat(nm); err == nil {
		goto gen
	}
	return nm
}

// IDToName resolves an ID to an actual name,
// with the help of a conf struct. If the name
// is not found, throws an error
func (c *Conf) IDToName(id string) (string, error) {
	nm := fmt.Sprintf("%s/%s", c.Prefix, id)
	if _, err := os.Stat(nm); err == nil {
		return nm, nil
	}

	return "", errors.New("file not found")
}
