package lib

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"lazygitRepoManager/internal/entity"
	"lazygitRepoManager/internal/vars"
)

type Configuration struct {
	Config            entity.ConfigurationEntity
	configurationFile string
}

func NewConfiguration() *Configuration {
	config := &Configuration{}
	config.loadConfigurations()
	return config
}

func (c *Configuration) getDefaultTerminalCommand() string {
	if platform.IsWindows() {
		return "start /MAX powershell.exe -Command " + vars.COMMAND_KEY
	} else if platform.IsLinux() {
		return "tilix --maximize -e \"" + vars.COMMAND_KEY + "\""
	}
	return ""
}

func (c *Configuration) getDefaultLazygitCommand() string {
	if platform.IsWindows() {
		return console.WhichIgnoreError("lazygit.exe")
	} else if platform.IsLinux() {
		return console.WhichIgnoreError("lazygit")
	}
	return ""
}

func (c *Configuration) loadConfigurations() {
	c.Config = entity.ConfigurationEntity{Repositories: []string{}, TerminalCommand: ""}
	c.configurationFile = file.JoinPath(system.HomeUserConfigDir(), vars.APP_NAME+".json")
	if file.IsFile(c.configurationFile) {
		data, err := file.ReadJsonFile[entity.ConfigurationEntity](c.configurationFile)
		logic.ProcessError(err)
		c.Config = data
	}
	if str.IsEmpty(c.Config.TerminalCommand) {
		c.Config.TerminalCommand = c.getDefaultTerminalCommand()
	}
	if str.IsEmpty(c.Config.LazygitCommand) {
		c.Config.LazygitCommand = c.getDefaultLazygitCommand()
	}
	c.Config.Repositories = slice.RemoveDuplicate(c.Config.Repositories)
	c.UpdateConfigurations()
}

/* -------------------------------------------------------------------------- */
/*                                 PUBLIC AREA                                */
/* -------------------------------------------------------------------------- */
func (c *Configuration) UpdateConfigurations() {
	file.WriteJsonFile(c.configurationFile, c.Config, false)
}

func (c *Configuration) HasRepositories() bool {
	return len(c.Config.Repositories) > 0
}

func (c *Configuration) ChangeTerminalCommand() {
	logger.Warn(fmt.Sprintf(`Keep %s, because will be replaced with lazygit command.`, vars.COMMAND_KEY))
	userInput := console.ReadBashUserInput("Insert the new command[PRESS ENTER TO CANCEL]")
	if !str.IsEmpty(userInput) {
		logger.Info(fmt.Sprintf(`Inserted: %s`, userInput))
		c.Config.TerminalCommand = userInput
		c.UpdateConfigurations()
	}
}

func (c *Configuration) ChangeLazygitCommand() {
	userInput := console.ReadBashUserInput("Insert the new command[PRESS ENTER TO CANCEL]")
	if !str.IsEmpty(userInput) {
		logger.Info(fmt.Sprintf(`Inserted: %s`, userInput))
		c.Config.LazygitCommand = userInput
		c.UpdateConfigurations()
	}
}
