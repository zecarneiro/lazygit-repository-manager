package main

import (
	"lazygitRepoManager/src/lib"

	"github.com/rivo/tview"
)

var (
	app *tview.Application
)

func main() {
	lib.Init()
	app = tview.NewApplication()
	lib.StartApp(app)
}
