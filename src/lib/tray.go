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
	configurations := fyne.NewMenuItem("Configurations", func() {
		openConfigurationsForm()
	})

	addNewRepository := fyne.NewMenuItem("Add new repository", func() {
		openAddNewRepoForm()
	})
	removeInvalidRepositories := fyne.NewMenuItem(REMOVE_IVALID_REPOSITORY_TRAY_MENU_TITLE, func() {
		removeInvalidRepositorories()
		reloadRepositoryMenuItem()
		Notify(REMOVE_IVALID_REPOSITORY_TRAY_MENU_TITLE + ", done.")
	})
	buildTrayRepositories(repositoriesMenuItem)
	trayMenu = fyne.NewMenu(APP_NAME)
	trayMenu.Items = append(trayMenu.Items, repositoriesMenuItem, addNewRepository, removeInvalidRepositories, fyne.NewMenuItemSeparator(), configurations)
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
	desk.SetSystemTrayIcon(fyne.NewStaticResource(APP_NAME, GetIcon()))
}
