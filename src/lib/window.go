package lib

import (
	"jnoronhautils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	formWindow          fyne.Window
	formWindowSize      fyne.Size
	directoryWindow     fyne.Window
	directoryWindowSize fyne.Size
)

func initWindow() {
	formWindowSize = fyne.NewSize(600, 500)
	directoryWindowSize = fyne.NewSize(600, 500)
}

// Directory Window and Dialog Area
func loadDirectoryWindow() {
	if directoryWindow != nil {
		directoryWindow.Close()
		directoryWindow = nil
	}
	directoryWindow = application.NewWindow(APP_NAME + " - Select Repository")
	directoryWindow.Resize(directoryWindowSize)
	directoryWindow.FixedSize()
	directoryWindow.CenterOnScreen()
}

func closeDirectoryWindow() {
	directoryWindow.Close()
	directoryWindow = nil
}

func openDialogFolder(callback func(fyne.ListableURI, error), uri fyne.ListableURI) {
	loadDirectoryWindow()
	dialogDir := dialog.NewFolderOpen(callback, directoryWindow)
	dialogDir.Resize(directoryWindowSize)
	dialogDir.SetOnClosed(closeDirectoryWindow)
	if uri != nil && jnoronhautils.FileExist(uri.Path()) {
		dialogDir.SetLocation(uri)
	}
	dialogDir.Show()
	directoryWindow.Show()
}

// Form Window and Dialog Area
func loadFormWindow() {
	if formWindow != nil {
		formWindow.Close()
		formWindow = nil
	}
	formWindow = application.NewWindow(APP_NAME + " - Configuration")
	formWindow.Resize(formWindowSize)
	formWindow.FixedSize()
	formWindow.CenterOnScreen()
}

func closeFormWindow() {
	formWindow.Close()
	formWindow = nil
}

func openDialogForm(callback func(status bool), title string, items []*widget.FormItem) {
	loadFormWindow()
	dialogForm := dialog.NewForm(title, "Save", "Cancel", items, callback, formWindow)
	dialogForm.Resize(directoryWindowSize)
	dialogForm.SetOnClosed(closeFormWindow)
	dialogForm.Show()
	formWindow.Show()
}
