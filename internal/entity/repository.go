package entity

import (
	"fmt"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/ui"
	"lazygitRepoManager/internal/shared"
	"slices"
)

type Repository struct {
	configuration *ConfigurationData
}

func NewRepository(configuration *ConfigurationData) *Repository {
	return &Repository{
		configuration: configuration,
	}
}

func (r *Repository) isValidGitRepository(repo string) bool {
	gitDir := file.JoinPath(repo, ".git")
	return file.IsDir(gitDir)
}

func (r *Repository) OpenRepository(repo string) {
	if !file.IsDir(repo) || !r.isValidGitRepository(repo) {
		ui.ErrorDialog(shared.AppDisplayName, fmt.Sprintf(`Invalid repository: "%s"`, repo))
	} else {
		cmdStr := fmt.Sprintf(`%s`, str.StringReplaceAll(r.configuration.Config.TerminalCommand, map[string]string{shared.COMMAND_KEY: r.configuration.Config.LazygitCommand}))
		cmdInfo := models.Command{
			Cmd:        cmdStr,
			Cwd:        repo,
			UseShell:   true,
			Verbose:    false,
			IsAsync:    true,
			Background: true,
		}
		shared.InfoNofity(fmt.Sprintf(`Openning repository: %s`, repo))
		logger.Error(exe.ExecRealTime(cmdInfo))
	}
}

func (r *Repository) AddNewRepo() {
	repoPathData := ui.SelectFolder(fmt.Sprintf(`%s - Select repository`, shared.AppDisplayName))
	if repoPathData.HasError() {
		shared.ErrorNofity(repoPathData.Error.Error())
	} else {
		if !str.IsEmpty(repoPathData.Data) {
			if !slices.Contains(r.configuration.Config.Repositories, repoPathData.Data) && r.isValidGitRepository(repoPathData.Data) {
				r.configuration.Config.Repositories = append(r.configuration.Config.Repositories, repoPathData.Data)
				r.configuration.UpdateConfigurations()
				shared.OkNofity(fmt.Sprintf(`Repo saved successfully: "%s"`, repoPathData.Data))
			} else {
				shared.ErrorNofity(fmt.Sprintf(`Invalid given Repo: "%s"`, repoPathData.Data))
			}
		}
	}
}

func (r *Repository) DeleteRepo() {
	if r.configuration.HasRepositories() {
		selectedData := ui.MultiSelectList(shared.AppDisplayName, "Select one or more repository to delete:", r.configuration.Config.Repositories)
		if selectedData.HasError() {
			shared.ErrorNofity(selectedData.Error.Error())
		} else {
			if len(selectedData.Data) > 0 {
				for _, repo := range selectedData.Data {
					r.configuration.Config.Repositories = slice.FilterArray(r.configuration.Config.Repositories, func(value string) bool {
						return value != repo
					})
				}
				r.configuration.UpdateConfigurations()
				shared.OkNofity("Repo(s) deleted")
			}
		}
	}
}

func (r *Repository) DeleteInvalidRepo() {
	shared.InfoNofity("Processing, delete invalid repositories...")
	validRepos := slice.FilterArray(r.configuration.Config.Repositories, func(value string) bool {
		return r.isValidGitRepository(value)
	})
	if len(validRepos) != len(r.configuration.Config.Repositories) {
		r.configuration.Config.Repositories = validRepos
		r.configuration.UpdateConfigurations()
	}
	shared.OkNofity("Processing, delete invalid repositories, Done.")
}
