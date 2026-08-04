package main

import (
	"golangutils/pkg/consolemenu"
	"lazygitRepoManager/internal/lib"
	"lazygitRepoManager/internal/vars"
)

var (
	app *consolemenu.ConsoleMenu
)

func main() {
	app = consolemenu.New()
	app.Title = vars.APP_DISPLAY_NAME
	lib.StartApp(app)
}
