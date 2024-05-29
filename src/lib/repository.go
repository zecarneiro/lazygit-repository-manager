package lib

import (
	"jnoronhautils"
	"jnoronhautils/entities"

	"fyne.io/fyne/v2"
)

var (
	lastOpenRepo fyne.ListableURI
)

func isValidGitRepository(repo string) bool {
	gitDir := jnoronhautils.ResolvePath(repo + "/.git")
	return jnoronhautils.FileExist(gitDir)
}

func buildTrayRepositories(menu *fyne.MenuItem) {
	menu.ChildMenu = fyne.NewMenu("Repositories Items")
	for _, repo := range config.Repositories {
		repo = jnoronhautils.ResolvePath(repo)
		title := repo
		if len(title) > config.MaxCharRepoRepresentation {
			title = jnoronhautils.GetSubstring(title, 0, config.MaxCharRepoRepresentation)
			title += "..."
			jnoronhautils.WarnLog("New repo representation: "+title, false)
		}
		if !jnoronhautils.FileExist(repo) {
			title = title + " - NOT FOUND"
		} else if !isValidGitRepository(repo) {
			title = title + " - INVALID GIT REPOSITORY"
		}
		menu.ChildMenu.Items = append(menu.ChildMenu.Items, fyne.NewMenuItem(title, func() {
			openRepository(repo)
		}))
	}
}

func openAddNewRepoForm() {
	openDialogFolder(ProcessNewRepo, lastOpenRepo)
}

func removeInvalidRepositorories() {
	newRepos := []string{}
	for _, repo := range config.Repositories {
		if jnoronhautils.FileExist(repo) && isValidGitRepository(repo) {
			newRepos = append(newRepos, repo)
		}
	}
	config.Repositories = newRepos
	updateConfigurations()
}

func openRepository(repo string) {
	if !jnoronhautils.FileExist(repo) || !isValidGitRepository(repo) {
		Notify("Invalid repository: " + repo + ". Please, run 'Remove invalid repositories'")
	} else {
		fullLazygitCmd := config.LazygitCommand + " -p '" + jnoronhautils.ResolvePath(repo) + "'"
		cmd := entities.CommandInfo{
			Cmd:           jnoronhautils.StringReplaceAll(config.TerminalCommand, map[string]string{COMMAND_KEY: fullLazygitCmd}),
			UsePowerShell: false,
		}
		jnoronhautils.InfoLog("Exec: "+cmd.Cmd, false)
		jnoronhautils.Exec(cmd)
	}
}

func ProcessNewRepo(uri fyne.ListableURI, err error) {
	if err != nil {
		jnoronhautils.ErrorLog(err.Error(), false)
	} else if uri != nil {
		lastOpenRepo = uri
		if !jnoronhautils.InArray(config.Repositories, uri.Path()) {
			config.Repositories = append(config.Repositories, uri.Path())
			updateConfigurations()
			reloadRepositoryMenuItem()
		}
	}
}
