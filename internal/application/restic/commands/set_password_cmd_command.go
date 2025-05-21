package commands

type SetPasswordCmdCommand struct {
	RootDir string
}

func (c SetPasswordCmdCommand) CommandName() string {
	return "SetPasswordCmd"
}
