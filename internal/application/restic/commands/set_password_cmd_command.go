package commands

type SetPasswordCmdCommand struct {
	RootDir  string
	Password string
}

func (c SetPasswordCmdCommand) CommandName() string {
	return "SetPasswordCmd"
}
