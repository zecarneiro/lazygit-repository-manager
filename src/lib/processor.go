package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zecarneiro/simpleconsoleui"

	"github.com/rivo/tview"
	"github.com/zecarneiro/golangutils"
)

const (
	COMMAND_KEY = "__COMMAND__"
)

var (
	app     *tview.Application
	windows []simpleconsoleui.Window
	// Application information variables
	Author                 string
	ApplicationReleaseDate string
	ApplicationInfoFile    string
	ApplicationName        string
	ApplicationDisplayName string
	ApplicationId          string
	ApplicationVersion     string
	ApplicationIcon        string

	// Others
	executableDir string
)

func getIcon() []byte {
	iconPath := golangutils.ResolvePath(executableDir + "/icon/" + ApplicationIcon)
	return golangutils.ReadFileInByte(iconPath)
}

func setExecutableDir(dir string) {
	executableDir = dir
}

func loadAppInformations(line string, err error) {
	golangutils.ProcessError(err)
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
	} else if strings.HasPrefix(line, "WIN_ICON") && golangutils.IsWindows() {
		_, after, _ := strings.Cut(line, "=")
		ApplicationIcon = after
	} else if strings.HasPrefix(line, "LINUX_ICON") && golangutils.IsLinux() {
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
		golangutils.ProcessError(errors.New(golangutils.INVALID_PLATFORM_MSG))
	}
}

/* -------------------------------------------------------------------------- */
/*                                 VIEWS AREA                                 */
/* -------------------------------------------------------------------------- */
func quit() {
	stop := func(canContinue bool) {
		if canContinue {
			app.Stop()
		}
	}
	simpleconsoleui.Confirm("Do you want to quit the application?", "Quit", "", stop)
}

func about() tview.Primitive {
	about := tview.NewTextView().SetText(fmt.Sprintf("Author: %s\nRelease Date: %s\nVersion: %s", Author, ApplicationReleaseDate, ApplicationVersion))
	return about
}

func showModalErr() {
	simpleconsoleui.Error("Error modal", "", nil)
}

func showModalInfo() {
	simpleconsoleui.Information("Information modal", "", nil)
}
func showModalOk() {
	simpleconsoleui.Ok("Ok modal", "", nil)
}
func showModalWarn() {
	simpleconsoleui.Warn("Warn modal", "", nil)
}

/* -------------------------------------------------------------------------- */
/*                                 PUBLIC AREA                                */
/* -------------------------------------------------------------------------- */
func Init() {
	executableDir = golangutils.GetExecutableDir()
	ApplicationInfoFile = golangutils.ResolvePath(executableDir + "/app-information")
	golangutils.ReadFileLineByLine(ApplicationInfoFile, loadAppInformations)
	validateStart()
	simpleconsoleui.InitUi(tview.Theme{})
}

func StartApp(application *tview.Application) {
	app = application
	loadConfigurations()
	windows = []simpleconsoleui.Window{
		{MenuName: "Repositories", MenuPage: respositories, HasLog: true},
		{MenuName: "Add new repository", MenuPage: addNewRepository, HasLog: false},
		{MenuName: "Remove invalid repositories", Callback: delInvalidRepositories},
		{MenuName: "Configurations", MenuPage: configuration, HasLog: true},
		{MenuName: "Refresh", MenuPage: nil, HasLog: false, Callback: simpleconsoleui.RefreshAndKeepOnPage},
		{MenuName: "About", MenuPage: about},
		{MenuName: "Quit", Callback: quit},
	}
	simpleconsoleui.Start(app, windows, ApplicationDisplayName, "Repositories manager")
}
