package internal

import (
	"golangutils/pkg/file"
	"lazygitRepoManager/internal/entity"
	"lazygitRepoManager/internal/shared"

	"github.com/energye/systray"
)

var (
	configuration    *entity.ConfigurationData
	repository       *entity.Repository
	isSystrayCreated bool
)

func showMenu(menu systray.IMenu) {
	if menu != nil {
		menu.ShowMenu()
	}
}

func buildSettingsMenu() {
	settingsMenuEntry := systray.AddMenuItem("Settings", "Settings")
	settingsMenuEntry.AddSubMenuItem("Add repository", "Add repository").Click(func() {
		repository.AddNewRepo()
		refresh(false)
	})
	deleteRepoMenuEntry := settingsMenuEntry.AddSubMenuItem("Delete repositories", "Delete repositories")
	if configuration.HasRepositories() {
		deleteRepoMenuEntry.Click(func() {
			repository.DeleteRepo()
			refresh(false)
		})
	} else {
		deleteRepoMenuEntry.Disable()
	}
	settingsMenuEntry.AddSubMenuItem("Delete invalid repositories", "Delete invalid repositories").Click(func() {
		repository.DeleteInvalidRepo()
		refresh(false)
	})
	settingsMenuEntry.AddSubMenuItem("Terminal Command", "Terminal Command").Click(func() {
		configuration.ChangeTerminalCommand()
	})
	settingsMenuEntry.AddSubMenuItem("Lazygit Command", "Lazygit Command").Click(func() {
		configuration.ChangeLazygitCommand()
	})
	buildThemeMenu(settingsMenuEntry)
}

func refresh(isStart bool) {
	if !isStart {
		systray.ResetMenu()
	}
	if configuration.HasRepositories() {
		for _, repo := range configuration.Config.Repositories {
			systray.AddMenuItem(repo, repo).Click(func() {
				repository.OpenRepository(repo)
				refresh(false)
			})
		}
		systray.AddSeparator()
	}
	buildSettingsMenu()
	systray.AddSeparator()
	systray.AddMenuItem("Exit", "Exit of the application").Click(func() {
		systray.Quit()
	})
}

func buildTrayApp() {
	refresh(true)
	if !isSystrayCreated {
		fileByte, _ := file.ReadFileInByte(shared.GetIcon())
		systray.SetIcon(fileByte)
		systray.SetTitle(shared.APP_DISPLAY_NAME)
		systray.SetTooltip(shared.APP_DISPLAY_NAME)
		systray.SetOnClick(showMenu)
		systray.SetOnRClick(showMenu)
		isSystrayCreated = true
	}
}

func Start() {
	configuration = entity.NewConfigurationData()
	repository = entity.NewRepository(configuration)
	isSystrayCreated = false
	systray.Run(buildTrayApp, nil)
}
