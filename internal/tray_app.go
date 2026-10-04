package internal

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/ui"
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

func buildGitConfigMenu(settingsMenu *systray.MenuItem) {
	usernameCmd := "git config --global user.name"
	emailCmd := "git config --global user.email"
	username, _ := exe.Exec(models.Command{Cmd: usernameCmd, UseShell: true, Background: true})
	email, _ := exe.Exec(models.Command{Cmd: emailCmd, UseShell: true, Background: true})
	usernameMenu := fmt.Sprintf(`Username global (Current: %s)`, username)
	emailMenu := fmt.Sprintf(`Email global (Current: %s)`, email)
	gitConfigMenuEntry := settingsMenu.AddSubMenuItem("Git user config", "Git user config")
	gitConfigMenuEntry.AddSubMenuItem(usernameMenu, usernameMenu).Click(func() {
		usernameData := ui.Input(shared.AppDisplayName, "Insert Git username:", username)
		if usernameData.HasError() {
			ui.ErrorNofity(usernameData.Error.Error(), shared.GetIcon())
		} else {
			err := exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(`%s "%s"`, usernameCmd, usernameData.Data), UseShell: true, Background: true})
			if err != nil {
				ui.ErrorNofity("Error on change git username", shared.GetIcon())
			} else {
				ui.OkNofity("Git username changed successfully.", shared.GetIcon())
				refresh(false)
			}
		}
	})
	gitConfigMenuEntry.AddSubMenuItem(emailMenu, emailMenu).Click(func() {
		emailData := ui.Input(shared.AppDisplayName, "Insert Git email:", email)
		if emailData.HasError() {
			ui.ErrorNofity(emailData.Error.Error(), shared.GetIcon())
		} else {
			err := exe.ExecRealTime(models.Command{Cmd: fmt.Sprintf(`%s "%s"`, emailCmd, emailData.Data), UseShell: true, Background: true})
			if err != nil {
				ui.ErrorNofity("Error on change git email", shared.GetIcon())
			} else {
				ui.OkNofity("Git email changed successfully.", shared.GetIcon())
				refresh(false)
			}
		}
	})
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
	buildGitConfigMenu(settingsMenuEntry)
	buildThemeMenu(settingsMenuEntry)
	settingsMenuEntry.AddSubMenuItem("About", "About").Click(func() {
		message := fmt.Sprintf("Name: %s", shared.APP_NAME)
		message += fmt.Sprintf("%sDisplay Name: %s", common.Eol(), shared.AppDisplayName)
		message += fmt.Sprintf("%sVersion: %s", common.Eol(), shared.AppVersion)
		message += fmt.Sprintf("%sRelease Date: %s", common.Eol(), shared.AppReleaseDate)
		logger.Error(ui.InfoDialog(shared.AppDisplayName, message))
	})
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
		systray.SetTitle(shared.AppDisplayName)
		systray.SetTooltip(shared.AppDisplayName)
		systray.SetOnClick(showMenu)
		systray.SetOnRClick(showMenu)
		isSystrayCreated = true
	}
}

func Start() {
	shared.LoadAppInformations()
	configuration = entity.NewConfigurationData()
	repository = entity.NewRepository(configuration)
	isSystrayCreated = false
	systray.Run(buildTrayApp, nil)
}
