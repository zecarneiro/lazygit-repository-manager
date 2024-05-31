package lib

import (
	"errors"
	"jnoronhautils"
	jnoronhautilsEnums "jnoronhautils/enums"
	"strings"

	"fyne.io/fyne/v2"
)

const (
	COMMAND_KEY = "__COMMAND__"
)

var (
	// Application information variables
	Author string
	ApplicationReleaseDate string
	ApplicationInfoFile string
	ApplicationName string
	ApplicationDisplayName string
	ApplicationId string
	ApplicationVersion string
	ApplicationIcon string

	// Others
	executableDir string
	application   fyne.App
)

func loadAppInformations(line string, err error) {
	jnoronhautils.ProcessError(err)
	if strings.HasPrefix(line, "NAME") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationName = after
	} else if strings.HasPrefix(line, "DISPLAY_NAME") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationDisplayName = after
	} else if strings.HasPrefix(line, "ID") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationId = after
	} else if strings.HasPrefix(line, "VERSION") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationVersion = after
	} else if strings.HasPrefix(line, "WIN_ICON") && jnoronhautils.IsWindows() {
		_, after, _ := strings.Cut(line, "=")
		ApplicationIcon = after
	} else if strings.HasPrefix(line, "LINUX_ICON") && jnoronhautils.IsLinux() {
		_, after, _ := strings.Cut(line, "=")
		ApplicationIcon = after
	} else if strings.HasPrefix(line, "AUTHOR") {
		_, after, _ := strings.Cut(line, "=")
		Author = after
	} else if strings.HasPrefix(line, "RELEASE_DATE") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationReleaseDate = after
	}
}

func validateStart() {
	if len(ApplicationIcon) == 0 {
		jnoronhautils.ProcessError(errors.New(jnoronhautilsEnums.INVALID_PLATFORM_MSG))
	}
}

func Init() {
	executableDir = jnoronhautils.GetExecutableDir()
	ApplicationInfoFile = jnoronhautils.ResolvePath(executableDir + "/app-information")
	jnoronhautils.ReadFileLineByLine(ApplicationInfoFile, loadAppInformations)
	validateStart()
}

func GetIcon() []byte {
	iconPath := jnoronhautils.ResolvePath(executableDir + "/icon/" + ApplicationIcon)
	return jnoronhautils.ReadFileInByte(iconPath)
}

func Notify(content string) {
	notification := fyne.NewNotification(ApplicationDisplayName, content)
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
