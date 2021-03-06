package common

// Conf encapsulates configuration details of the
// ethersim software for this run.
type Conf struct {
	MinID   int
	MaxID   int
	Prefix  string
	Serving bool
	Name    string
}

// MkConfig generates a configuration
// for this run. Currently just uses default
// values.
func MkConfig() Conf {
	return Conf{
		MinID:   defMinID,
		MaxID:   defMaxID,
		Prefix:  defPrefix,
		Serving: false,
	}
}
