package entity

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
	"golangutils/pkg/ui"
	"lazygitRepoManager/internal/shared"
)

type ConfigurationData struct {
	Config            Configuration
	configurationFile string
}

func NewConfigurationData() *ConfigurationData {
	config := &ConfigurationData{}
	config.loadConfigurations()
	return config
}

func (c *ConfigurationData) getDefaultTerminalCommand() string {
	if platform.IsWindows() {
		return fmt.Sprintf(`wt.exe -M powershell -Command "%s"`, shared.COMMAND_KEY)
	} else if platform.IsLinux() {
		return fmt.Sprintf(`ghostty --maximize -e "%s"`, shared.COMMAND_KEY)
	}
	return ""
}

func (c *ConfigurationData) getDefaultLazygitCommand() string {
	if platform.IsWindows() {
		return console.WhichIgnoreError("lazygit.exe")
	} else if platform.IsLinux() {
		return console.WhichIgnoreError("lazygit")
	}
	return ""
}

func (c *ConfigurationData) loadConfigurations() {
	c.Config = Configuration{Repositories: []string{}, TerminalCommand: ""}
	c.configurationFile = file.JoinPath(system.HomeUserConfigDir(), shared.APP_NAME+".json")
	if file.IsFile(c.configurationFile) {
		data, err := file.ReadJsonFile[Configuration](c.configurationFile)
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
func (c *ConfigurationData) UpdateConfigurations() {
	file.WriteJsonFile(c.configurationFile, c.Config, false)
}

func (c *ConfigurationData) HasRepositories() bool {
	return len(c.Config.Repositories) > 0
}

func (c *ConfigurationData) ChangeTerminalCommand() {
	bodyData := fmt.Sprintf(`Keep %s, because will be replaced with lazygit command.`, shared.COMMAND_KEY)
	title := fmt.Sprintf(`%s - Insert the new terminal command`, shared.AppDisplayName)
	responseData := ui.Input(title, bodyData, c.Config.TerminalCommand)
	if responseData.HasError() {
		ui.ErrorNofity(responseData.Error.Error(), shared.GetIcon())
	} else {
		c.Config.TerminalCommand = responseData.Data
		c.UpdateConfigurations()
	}
}

func (c *ConfigurationData) ChangeLazygitCommand() {
	title := fmt.Sprintf(`%s - Insert the Lazygit command`, shared.AppDisplayName)
	responseData := ui.Input(title, "Insert full path or the command itself", c.Config.LazygitCommand)
	if responseData.HasError() {
		ui.ErrorNofity(responseData.Error.Error(), shared.GetIcon())
	} else {
		c.Config.LazygitCommand = responseData.Data
		c.UpdateConfigurations()
	}
}
