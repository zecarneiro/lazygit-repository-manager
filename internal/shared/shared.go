package shared

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
	"strings"
)

const (
	CONFIG_RESET_MESSAGE = "To reset, keep fields empty"
	APP_NAME             = "lazygit-repository-manager"
	COMMAND_KEY          = "__COMMAND__"
)

var (
	AppDisplayName string
	AppVersion     string
	AppReleaseDate string
)

func loadAppInformations(line string) {
	if strings.HasPrefix(line, "DISPLAY_NAME") {
		_, after, _ := strings.Cut(line, "=")
		AppDisplayName = after
	} else if strings.HasPrefix(line, "VERSION") {
		_, after, _ := strings.Cut(line, "=")
		AppVersion = after
	} else if strings.HasPrefix(line, "RELEASE_DATE") {
		_, after, _ := strings.Cut(line, "=")
		AppReleaseDate = after
	}
}

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

func LoadAppInformations() {
	logger.Error(file.ReadFileLineByLine(file.JoinPath(GetExecutableDir(), "APP_INFO.conf"), loadAppInformations))
}
