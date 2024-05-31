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
	informationWindow     fyne.Window
	informationWindowSize     fyne.Size
)

func initWindow() {
	formWindowSize = fyne.NewSize(600, 500)
	directoryWindowSize = fyne.NewSize(600, 500)
	informationWindowSize = fyne.NewSize(300, 200)
}

// Directory Window and Dialog Area
func loadDirectoryWindow() {
	if directoryWindow != nil {
		directoryWindow.Close()
		directoryWindow = nil
	}
	directoryWindow = application.NewWindow(ApplicationDisplayName + " - Select Repository")
	directoryWindow.Resize(directoryWindowSize)
	directoryWindow.SetFixedSize(true)
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
	formWindow = application.NewWindow(ApplicationDisplayName + " - Configuration")
	formWindow.Resize(formWindowSize)
	formWindow.SetFixedSize(true)
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

func loadInformationWindow() {
	if informationWindow != nil {
		informationWindow.Close()
		informationWindow = nil
	}
	informationWindow = application.NewWindow(ApplicationDisplayName + " - Select Repository")
	informationWindow.Resize(informationWindowSize)
	informationWindow.SetFixedSize(true)
	informationWindow.CenterOnScreen()
}

func closeInformationWindow() {
	informationWindow.Close()
	informationWindow = nil
}

func openInformation(title string, message string) {
	loadInformationWindow()
	dialogInfo := dialog.NewInformation(title, message, informationWindow)
	dialogInfo.SetOnClosed(closeInformationWindow)
	dialogInfo.Resize(informationWindowSize)
	dialogInfo.Show()
	informationWindow.Show()
}
