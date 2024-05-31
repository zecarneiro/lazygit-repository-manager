package lib

import (
	"jnoronhautils"
	"lazygitRepoManager/src/entity"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/widget"
)

const (
	CONFIG_RESET_MESSAGE = "To reset, keep fields empty"
)

var (
	config            entity.Configuration
	configurationFile string

	// Configurations Form items
	inputTerminalCommand           *widget.Entry
	inputLazygitCommand            *widget.Entry
	inputMaxCharRepoRepresentation *widget.Entry
)

func loadMaxCharRepoRepresentation() {
	if config.MaxCharRepoRepresentation == 0 {
		config.MaxCharRepoRepresentation = 100
	}
}

func loadTerminalCommand() {
	if len(config.TerminalCommand) == 0 {
		if jnoronhautils.IsWindows() {
			config.TerminalCommand = "start /MAX powershell.exe -Command " + COMMAND_KEY
		}
	}
}

func loadLazygitCommand() {
	if len(config.LazygitCommand) == 0 {
		if jnoronhautils.IsWindows() {
			config.LazygitCommand = "lazygit.exe"
		}
	}
}

func loadConfigurations() {
	config = entity.Configuration{Repositories: []string{}, TerminalCommand: ""}

	// Set Configuration Dir
	configurationDir := jnoronhautils.ResolvePath(jnoronhautils.SystemInfo().HomeDir + "/.config/")
	jnoronhautils.CreateDirectory(configurationDir, true)

	configurationFile = jnoronhautils.ResolvePath(configurationDir + "/" + ApplicationName + ".json")
	if jnoronhautils.FileExist(configurationFile) {
		data, err := jnoronhautils.ReadJsonFile[entity.Configuration](configurationFile)
		if err != nil {
			jnoronhautils.ErrorLog(err.Error(), false)
		} else {
			config = data
		}
	}
	loadTerminalCommand()
	loadLazygitCommand()
	loadMaxCharRepoRepresentation()
	config.Repositories = jnoronhautils.RemoveDuplicate(config.Repositories)
	updateConfigurations()
}

func updateConfigurations() {
	jnoronhautils.WriteJsonFile(configurationFile, config)
}

func getTerminalCommandFormItem() *widget.FormItem {
	inputTerminalCommand = widget.NewEntry()
	inputTerminalCommand.Text = config.TerminalCommand
	inputTerminalCommand.Show()
	return widget.NewFormItem("Terminal Command", inputTerminalCommand)
}

func getMaxCharRepoRepresntationFormItem() *widget.FormItem {
	inputMaxCharRepoRepresentation = widget.NewEntry()
	inputMaxCharRepoRepresentation.Text = strconv.Itoa(config.MaxCharRepoRepresentation)
	inputMaxCharRepoRepresentation.Show()
	return widget.NewFormItem("Max char repo representation", inputMaxCharRepoRepresentation)
}

func getLazygitCommandFormItem() *widget.FormItem {
	inputLazygitCommand = widget.NewEntry()
	inputLazygitCommand.Text = config.LazygitCommand
	inputLazygitCommand.Show()
	return widget.NewFormItem("Lazygit Command", inputLazygitCommand)
}

func openConfigurationsForm() {
	items := []*widget.FormItem{getLazygitCommandFormItem(), getTerminalCommandFormItem(), getMaxCharRepoRepresntationFormItem()}
	openDialogForm(ProcessConfiguration, CONFIG_RESET_MESSAGE, items)
}

func ProcessConfiguration(status bool) {
	if status {
		if len(inputTerminalCommand.Text) > 0 && !strings.Contains(inputTerminalCommand.Text, COMMAND_KEY) {
			Notify("Terminal Command must contain: " + COMMAND_KEY)
		} else {
			config.LazygitCommand = inputLazygitCommand.Text
			loadLazygitCommand()
			config.TerminalCommand = inputTerminalCommand.Text
			loadTerminalCommand()
			value, err := strconv.Atoi(inputMaxCharRepoRepresentation.Text)
			jnoronhautils.HasAndLogError(err)
			oldMaxCharRepoRepresentation := config.MaxCharRepoRepresentation
			config.MaxCharRepoRepresentation = value
			loadMaxCharRepoRepresentation()
			// Save configuration
			updateConfigurations()
			if oldMaxCharRepoRepresentation != value {
				reloadRepositoryMenuItem()
			}
		}
	}
}
