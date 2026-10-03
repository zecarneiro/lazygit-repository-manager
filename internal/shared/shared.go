package shared

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
)

const (
	CONFIG_RESET_MESSAGE = "To reset, keep fields empty"
	APP_NAME             = "lazygit-repository-manager"
	APP_DISPLAY_NAME     = "Lazygit Repository Manager"
	COMMAND_KEY          = "__COMMAND__"
)

func GetExecutableDir() string {
	if platform.IsLinux() {
		return file.JoinPath(system.HomeUserOptDir(), APP_NAME)
	}
	directory, _ := exe.GetExecutableDir()
	return directory
}

func GetIcon() string {
	icon := ""
	if platform.IsWindows() {
		icon = fmt.Sprintf(`%s/win.ico`, GetExecutableDir())
	} else if platform.IsLinux() {
		icon = fmt.Sprintf(`%s/linux.png`, GetExecutableDir())
	}
	return file.ResolvePath(icon)
}
