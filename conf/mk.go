package conf

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Conf encapsulates configuration details of the
// ethersim software for this run.
type Conf struct {
	MinID  int
	MaxID  int
	Prefix string
}

// MkConfig generates a configuration
// for this run. Currently just uses default
// values.
func MkConfig() (Conf, error) {
	return Conf{
		MinID:  defMinID,
		MaxID:  defMaxID,
		Prefix: defPrefix,
	}, nil
}

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
