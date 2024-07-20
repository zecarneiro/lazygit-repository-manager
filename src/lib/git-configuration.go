package lib

import (
	"strings"

	"github.com/rivo/tview"
	"github.com/zecarneiro/golangutils"
	"github.com/zecarneiro/simpleconsoleui"
)

func getGitCmd(args []string) golangutils.CommandInfo {
	cmd := golangutils.CommandInfo{
		Cmd:  "git",
		Args: []string{"config", "--global"},
	}
	cmd.Args = append(cmd.Args, args...)
	if golangutils.IsWindows() {
		cmd.UsePowerShell = false
	} else if golangutils.IsLinux() {
		cmd.UseBash = true
	}
	return cmd
}

func getGitCredentialsCmd(isSet bool) golangutils.CommandInfo {
	cmd := getGitCmd([]string{"credential.helper"})
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

func getGitEmail() string {
	cmd := getGitCmd([]string{"user.email"})
	simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
	resp := golangutils.Exec(cmd)
	return strings.Trim(resp.Data, golangutils.SysInfo().Eol)
}
func setGitEmail(email string) {
	cmd := getGitCmd([]string{"user.email", "'" + email + "'"})
	simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
	golangutils.Exec(cmd)
}

func getGitName() string {
	cmd := getGitCmd([]string{"user.name"})
	simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
	resp := golangutils.Exec(cmd)
	return strings.Trim(resp.Data, golangutils.SysInfo().Eol)
}
func setGitName(name string) {
	cmd := getGitCmd([]string{"user.name", "'" + name + "'"})
	simpleconsoleui.PromptLog(golangutils.GetCommandToRun(golangutils.AddShellCommand(cmd)))
	golangutils.Exec(cmd)
}

/* -------------------------------------------------------------------------- */
/*                                 VIEWS AREA                                 */
/* -------------------------------------------------------------------------- */
func gitConfiguration() tview.Primitive {
	gitCredentialsStatus := getGitCredentialsStatus()
	email := getGitEmail()
	name := getGitName()
	formConfig := tview.NewForm()
	formConfig.AddTextView("Information", "All those configuration will be set as Global", 0, 2, true, true)
	formConfig.AddInputField("Email", email, 0, nil, func(text string) {
		email = text
	})
	formConfig.AddInputField("Name", name, 0, nil, func(text string) {
		name = text
	})
	formConfig.AddCheckbox("Enable Store Git Credentials", gitCredentialsStatus, func(checked bool) {
		gitCredentialsStatus = checked
	})
	formConfig.AddButton("Save", func() {
		setGitEmail(email)
		setGitName(name)
		if gitCredentialsStatus {
			enableGitCredentials()
		} else {
			disableGitCredentials()
		}
		simpleconsoleui.Ok("Configuration saved successfully", "", nil)
		formConfig.GetButton(0).Blur()
	})
	return formConfig
}
