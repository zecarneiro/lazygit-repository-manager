package main

import (
	"lazygitRepoManager/src/lib"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	application fyne.App
)

func main() {
	lib.Init()
	application = app.NewWithID(lib.ApplicationId)
	application.SetIcon(fyne.NewStaticResource(lib.ApplicationDisplayName, lib.GetIcon()))
	lib.StartApp(application)
	application.Run()
}
