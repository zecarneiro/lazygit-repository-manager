package lib

import (
	"errors"
	"jnoronhautils"
	jnoronhautilsEnums "jnoronhautils/enums"

	"fyne.io/fyne/v2"
)

const (
	COMMAND_KEY = "__COMMAND__"
	APP_NAME    = "Lazygit Repo Manager"
	APP_ID      = "lazygit-repo-manager"
)

var (
	executableDir string
	application   fyne.App
)

func GetIcon() []byte {
	iconPath := jnoronhautils.ResolvePath(executableDir + "/icon")
	if jnoronhautils.IsWindows() {
		iconPath = jnoronhautils.ResolvePath(iconPath + "/win-icon.ico")
	} else if jnoronhautils.IsLinux() {
		iconPath = jnoronhautils.ResolvePath(iconPath + "/linux-icon.png")
	} else {
		jnoronhautils.ProcessError(errors.New(jnoronhautilsEnums.INVALID_PLATFORM_MSG))
	}
	return jnoronhautils.ReadFileInByte(string(iconPath))
}

func Notify(content string) {
	notification := fyne.NewNotification(APP_NAME, content)
	application.SendNotification(notification)
}

func SetExecutableDir(dir string) {
	executableDir = dir
}

func StartApp(fyneApplication fyne.App) {
	application = fyneApplication
	loadConfigurations()
	initWindow()
	initTray()
}
