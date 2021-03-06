package common

// ID RNG
const defMinID int = 1000000000
const defMaxID int = 9999999999

// MsgInterval denotes how often are messages sent, in ms
// TODO: store this in Conf
const MsgInterval int = 100

// Where files of the form .../$id
// are stored, where $id is between MinID and MaxID
const defPrefix string = "/tmp"
