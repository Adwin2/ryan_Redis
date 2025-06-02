package command

func isWriteCmd(cmd string) bool {
	if cmd == "SET" {
		return true
	}
	return false
}
