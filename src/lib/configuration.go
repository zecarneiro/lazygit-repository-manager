package lib

import (
	"lazygitRepoManager/src/entity"
	"strings"

	"github.com/zecarneiro/simpleconsoleui"

	"github.com/rivo/tview"
	"github.com/zecarneiro/golangutils"
)

const (
	CONFIG_RESET_MESSAGE = "To reset, keep fields empty"
)

var (
	config            entity.Configuration
	configurationFile string
)

func getDefaultTerminalCommand() string {
	if golangutils.IsWindows() {
		return "start /MAX powershell.exe -Command " + COMMAND_KEY
	} else if golangutils.IsLinux() {
		return "tilix --maximize -e \"" + COMMAND_KEY + "\""
	}
	return ""
}
func getDefaultLazygitCommand() string {
	if golangutils.IsWindows() {
		return "lazygit.exe"
	} else if golangutils.IsLinux() {
		return "lazygit"
	}
	return ""
}

func getGitCredentialsCmd(isSet bool) golangutils.CommandInfo {
	cmd := golangutils.CommandInfo{
		Cmd:  "git",
		Args: []string{"config", "--global", "credential.helper"},
	}
	if golangutils.IsWindows() {
		cmd.UsePowerShell = false
	} else if golangutils.IsLinux() {
		cmd.UseBash = true
	}
	if isSet {
		cmd.Args = append(cmd.Args, "store")
	}
	return cmd
}
func getGitCredentialsStatus() bool {
	cmd := getGitCredentialsCmd(false)
	simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
	resp := golangutils.Exec(cmd)
	return strings.Contains(resp.Data, "store")
}
func enableGitCredentials() {
	cmd := getGitCredentialsCmd(true)
	simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
	golangutils.ExecRealTime(cmd)
}
func disableGitCredentials() {
	fileToDelete := golangutils.ResolvePath(golangutils.SysInfo().HomeDir + "/.git-credentials")
	cmd := getGitCredentialsCmd(false)
	cmd.Args = append(cmd.Args, "''")
	simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
	golangutils.ExecRealTime(cmd)
	simpleconsoleui.InfoLog("Delete file: " + fileToDelete)
	golangutils.DeleteFile(fileToDelete)
}

func loadConfigurations() {
	config = entity.Configuration{Repositories: []string{}, TerminalCommand: ""}

	// Set Configuration Dir
	configurationDir := golangutils.ResolvePath(golangutils.SysInfo().HomeDir + "/.config/")
	golangutils.CreateDirectory(configurationDir, true)

	configurationFile = golangutils.ResolvePath(configurationDir + "/" + ApplicationName + ".json")
	if golangutils.FileExist(configurationFile) {
		data, err := golangutils.ReadJsonFile[entity.Configuration](configurationFile)
		if err != nil {
			golangutils.ErrorLog(err.Error(), false)
		} else {
			config = data
		}
	}
	if len(config.TerminalCommand) == 0 {
		config.TerminalCommand = getDefaultTerminalCommand()
	}
	if len(config.LazygitCommand) == 0 {
		config.LazygitCommand = getDefaultLazygitCommand()
	}
	config.Repositories = golangutils.RemoveDuplicate(config.Repositories)
	updateConfigurations()
}
func updateConfigurations() {
	golangutils.WriteJsonFile(configurationFile, config)
}

/* -------------------------------------------------------------------------- */
/*                            VALIDATION FUNCS AREA                           */
/* -------------------------------------------------------------------------- */
func validateEmptyField(field string) bool {
	return len(field) > 0
}
func validateFieldContainsCommandKey(field string) bool {
	return strings.Contains(field, COMMAND_KEY)
}
func validateLazygitCommandField(value string) bool {
	if validateEmptyField(value) {
		return true
	}
	simpleconsoleui.ErrorLog("Invalid Lazygit Command")
	return false
}
func validateTerminalCommandField(value string) bool {
	if validateEmptyField(value) && validateFieldContainsCommandKey(value) {
		return true
	}
	simpleconsoleui.ErrorLog("Invalid Terminal Command")
	return false
}

/* -------------------------------------------------------------------------- */
/*                                 VIEWS AREA                                 */
/* -------------------------------------------------------------------------- */
func configuration() tview.Primitive {
	gitCredentialsStatus := getGitCredentialsStatus()
	formConfig := tview.NewForm()
	addField := func() {
		formConfig.AddInputField("Lazygit Command", config.LazygitCommand, 0, nil, func(text string) {
			config.LazygitCommand = text
		})
		formConfig.AddInputField("Terminal Command", config.TerminalCommand, 0, nil, func(text string) {
			config.TerminalCommand = text
		})
		formConfig.AddCheckbox("Enable Store Git Credentials", gitCredentialsStatus, func(checked bool) {
			gitCredentialsStatus = checked
		})
	}
	addField()
	formConfig.AddButton("Save", func() {
		if validateLazygitCommandField(config.LazygitCommand) && validateTerminalCommandField(config.TerminalCommand) {
			updateConfigurations()
			simpleconsoleui.Ok("Configuration saved successfully", "", nil)
			simpleconsoleui.ClearLog()
		}
		if gitCredentialsStatus {
			enableGitCredentials()
		} else {
			disableGitCredentials()
		}
		formConfig.GetButton(0).Blur()
	})
	formConfig.AddButton("Set default values", func() {
		config.LazygitCommand = getDefaultLazygitCommand()
		config.TerminalCommand = getDefaultTerminalCommand()
		formConfig.Clear(false)
		addField()
		formConfig.GetButton(1).Blur()
	})
	return formConfig
}
