package lib

import (
	"github.com/zecarneiro/simpleconsoleui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zecarneiro/golangutils"
)

var (
	repos *tview.List
)

func isValidGitRepository(repo string) bool {
	gitDir := golangutils.ResolvePath(repo + "/.git")
	return golangutils.FileExist(gitDir)
}

func processNewRepo(path string) {
	if len(path) > 0 {
		if !golangutils.InArray(config.Repositories, path) && isValidGitRepository(path) {
			config.Repositories = append(config.Repositories, path)
			updateConfigurations()
			reloadRepos()
			simpleconsoleui.Ok("Repo saved successfully: "+path, "", nil)
		} else {
			simpleconsoleui.Error("Invalid given Repo: "+path, "", nil)
		}
	}
}

func removeInvalidRepositorories() []string {
	invalidRepos := []string{}
	newRepos := []string{}
	for _, repo := range config.Repositories {
		if golangutils.FileExist(repo) && isValidGitRepository(repo) {
			newRepos = append(newRepos, repo)
		} else {
			invalidRepos = append(invalidRepos, repo)
		}
	}
	config.Repositories = newRepos
	updateConfigurations()
	reloadRepos()
	return invalidRepos
}

func openRepository(repo string) {
	if !golangutils.FileExist(repo) || !isValidGitRepository(repo) {
		simpleconsoleui.Error("Invalid repository: "+repo+". Please run 'Remove invalid repositories'", "", nil)
	} else {
		fullLazygitCmd := config.LazygitCommand + " -p '" + golangutils.ResolvePath(repo) + "'"
		cmd := golangutils.CommandInfo{
			Cmd: golangutils.StringReplaceAll(config.TerminalCommand, map[string]string{COMMAND_KEY: fullLazygitCmd}),
		}
		if golangutils.IsWindows() {
			cmd.UsePowerShell = false
		} else if golangutils.IsLinux() {
			cmd.UseBash = true
		}
		simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
		golangutils.ExecRealTimeAsync(cmd)
	}
}

/* -------------------------------------------------------------------------- */
/*                                 VIEWS AREA                                 */
/* -------------------------------------------------------------------------- */
func reloadRepos() {
	repos.Clear()
	repos.ShowSecondaryText(false)
	addRepos()
}
func addRepos() {
	for _, repo := range config.Repositories {
		repos.AddItem(repo, "", 0, nil)
	}
}
func respositories() tview.Primitive {
	repos = tview.NewList()
	open := func() {
		repo, _ := repos.GetItemText(repos.GetCurrentItem())
		openRepository(repo)
	}
	deleteRepo := func()  {
		repoSelected, _ := repos.GetItemText(repos.GetCurrentItem())
		simpleconsoleui.Confirm("Will remove the repository: " + repoSelected, "", "", func(canContinue bool) {
			if canContinue {
				config.Repositories = golangutils.FilterArray(config.Repositories, func(repo string) bool {
					return repo != repoSelected
				})
				updateConfigurations()
				reloadRepos()
			}
		})
	}
	repos.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			open()
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'd' {
			deleteRepo()
		}
		return event
	})
	repos.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftDoubleClick {
			open()
		}
		return action, event
	})
	reloadRepos()
	return repos
}
func addNewRepository() tview.Primitive {
	return simpleconsoleui.SelectTreeView(golangutils.SysInfo().HomeDir, true, false, "", processNewRepo)
}
func delInvalidRepositories() {
	simpleconsoleui.Confirm("Will remove all invalid repositories", "", "", func(canContinue bool) {
		if canContinue {
			for _, repo := range removeInvalidRepositorories() {
				simpleconsoleui.OkLog("Removed invalid repository: " + repo)
			}
		}
	})
}
