package command

// Command is the interface for the commands
type Command interface {
	Execute() (string, error)
}
