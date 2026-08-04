package lib

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/conv"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"lazygitRepoManager/internal/vars"
	"slices"
)

type Repository struct {
	configuration *Configuration
}

func NewRepository(configuration *Configuration) *Repository {
	return &Repository{
		configuration: configuration,
	}
}

func (r *Repository) isValidGitRepository(repo string) bool {
	gitDir := file.JoinPath(repo, ".git")
	return file.IsDir(gitDir)
}

func (r *Repository) openRepository(repo string) {
	if !file.IsDir(repo) || !r.isValidGitRepository(repo) {
		logger.Error(fmt.Errorf(`Invalid repository: "%s". %sPlease run 'Remove invalid repositories'`, repo, common.Eol()))
	} else {
		logger.Info(fmt.Sprintf(`Openning repository: %s`, repo))
		fullLazygitCmd := fmt.Sprintf(`%s`, r.configuration.Config.LazygitCommand)
		cmdInfo := models.Command{
			Cmd:      str.StringReplaceAll(r.configuration.Config.TerminalCommand, map[string]string{vars.COMMAND_KEY: fullLazygitCmd}),
			Cwd:      repo,
			UseShell: true,
			Verbose:  false,
			IsAsync:  true,
		}
		exe.Exec(cmdInfo)
	}
}

func (r *Repository) addNewRepo() {
	repoPath := console.ReadBashUserInput(fmt.Sprintf(`%s (PRESS ENTER TO IGNORE)`, "Insert path of repository"))
	if !str.IsEmpty(repoPath) {
		if !slices.Contains(r.configuration.Config.Repositories, repoPath) && r.isValidGitRepository(repoPath) {
			r.configuration.Config.Repositories = append(r.configuration.Config.Repositories, repoPath)
			r.configuration.UpdateConfigurations()
			logger.Ok(fmt.Sprintf(`Repo saved successfully: "%s"`, repoPath))
		} else {
			logger.Error(fmt.Errorf(`Invalid given Repo: "%s"`, repoPath))
		}
	}
}

func (r *Repository) deleteRepo() {
	if configuration.HasRepositories() {
		selectedRepo := -1
		minIndex := 1
		maxIndex := len(configuration.Config.Repositories)
		for {
			userInput := console.ReadBashUserInput(fmt.Sprintf(`Insert repository to delete(%d ... %d)[PRESS ENTER TO CANCEL]: `, minIndex, maxIndex))
			if str.IsEmpty(userInput) {
				break
			}
			userInputInt, err := conv.StringToInt(userInput)
			if err != nil || userInputInt < minIndex || userInputInt > maxIndex {
				logger.Warn("Insert a valid repository")
			} else {
				selectedRepo = userInputInt
				break
			}
		}
		if selectedRepo > 0 {
			repo := configuration.Config.Repositories[selectedRepo-1]
			if console.Confirm(fmt.Sprintf(`Will remove the repository: %s. Continue`, repo), true) {
				configuration.Config.Repositories = slice.FilterArray(configuration.Config.Repositories, func(value string) bool {
					return value != repo
				})
			}
			configuration.UpdateConfigurations()
		}
	}
}

func (r *Repository) deleteInvalidRepo() {
	validRepos := slice.FilterArray(configuration.Config.Repositories, func(value string) bool {
		return r.isValidGitRepository(value)
	})
	if len(validRepos) != len(configuration.Config.Repositories) {
		configuration.Config.Repositories = validRepos
		configuration.UpdateConfigurations()
	}
	logger.Ok("Done.")
}
