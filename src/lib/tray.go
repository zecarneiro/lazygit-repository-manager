package lib

import (
	"errors"
	"jnoronhautils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

const (
	REMOVE_IVALID_REPOSITORY_TRAY_MENU_TITLE = "Remove invalid repositories"
)

var (
	trayMenu             *fyne.Menu
	desk                 desktop.App
	repositoriesMenuItem *fyne.MenuItem
)

func buildTrayMenu() {
	repositoriesMenuItem = fyne.NewMenuItem("Repositories", func() {})
	addNewRepository := fyne.NewMenuItem("Add new repository", func() {
		openAddNewRepoForm()
	})
	removeInvalidRepositories := fyne.NewMenuItem(REMOVE_IVALID_REPOSITORY_TRAY_MENU_TITLE, func() {
		removeInvalidRepositorories()
		reloadRepositoryMenuItem()
		Notify(REMOVE_IVALID_REPOSITORY_TRAY_MENU_TITLE + ", done.")
	})
	buildTrayRepositories(repositoriesMenuItem)

	configurations := fyne.NewMenuItem("Configurations", func() {
		openConfigurationsForm()
	})
	about := fyne.NewMenuItem("About", func() {
		data := "Author: " + Author
		data += "\nRelease Date: " + ApplicationReleaseDate
		data += "\nVersion: " + ApplicationVersion
		openInformation(ApplicationDisplayName, data)
	})

	trayMenu = fyne.NewMenu(ApplicationDisplayName)
	trayMenu.Items = append(trayMenu.Items, repositoriesMenuItem, addNewRepository, removeInvalidRepositories, fyne.NewMenuItemSeparator(), configurations, about)
}

func reloadRepositoryMenuItem() {
	buildTrayRepositories(repositoriesMenuItem)
	trayMenu.Refresh()
}

func initTray() {
	appDesk, ok := application.(desktop.App)
	if !ok {
		jnoronhautils.ProcessError(errors.New("Error loading tray app"))
	}
	desk = appDesk
	buildTrayMenu()
	desk.SetSystemTrayMenu(trayMenu)
	desk.SetSystemTrayIcon(fyne.NewStaticResource(ApplicationDisplayName, GetIcon()))
}
