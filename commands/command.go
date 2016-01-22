package commands

//Command is a toplevel command to be executed by the cli's main routine
type Command interface {
	//Execute the command
	Execute() error
}