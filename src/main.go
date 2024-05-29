package main

import (
	"jnoronhautils"
	"lazygitRepoManager/src/lib"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var (
	application fyne.App
)

func main() {
	application = app.NewWithID(lib.APP_ID)
	lib.SetExecutableDir(jnoronhautils.GetExecutableDir())
	application.SetIcon(fyne.NewStaticResource(lib.APP_NAME, lib.GetIcon()))
	lib.StartApp(application)
	application.Run()
}
