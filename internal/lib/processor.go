package lib

import (
	"fmt"
	"golangutils/pkg/console"
	"golangutils/pkg/consolemenu"
)

var (
	configuration *Configuration
	repository    *Repository
)

func buildMenu(app *consolemenu.ConsoleMenu) {
	app.Reset()
	app.PreBuildCaller = func() {
		console.Clear()
	}
	app.PostProcessorCaller = func() {
		console.Pause()
		buildMenu(app)
	}
	if configuration.HasRepositories() {
		for _, repo := range configuration.Config.Repositories {
			app.AddEntry(repo, func(entry consolemenu.ConsoleMenuEntry) {
				repository.openRepository(entry.Label)
			})
		}
		app.AddSeparator()
	}
	app.AddEntry("Add new repository", func(entry consolemenu.ConsoleMenuEntry) {
		repository.addNewRepo()
	})
	if configuration.HasRepositories() {
		app.AddEntry("Delete repository", func(entry consolemenu.ConsoleMenuEntry) {
			repository.deleteRepo()
		})
		app.AddEntry("Delete invalid repositories", func(entry consolemenu.ConsoleMenuEntry) {
			repository.deleteInvalidRepo()
		})
	}
	app.AddSeparator()
	app.AddEntry(fmt.Sprintf(`Change Terminal Command. (Current: %s)`, configuration.Config.TerminalCommand), func(entry consolemenu.ConsoleMenuEntry) {
		configuration.ChangeTerminalCommand()
	})
	app.AddEntry(fmt.Sprintf(`Change Lazygit Command. (Current: %s)`, configuration.Config.LazygitCommand), func(entry consolemenu.ConsoleMenuEntry) {
		configuration.ChangeLazygitCommand()
	})
	app.Start()
}

/* -------------------------------------------------------------------------- */
/*                                 PUBLIC AREA                                */
/* -------------------------------------------------------------------------- */
func StartApp(app *consolemenu.ConsoleMenu) {
	configuration = NewConfiguration()
	repository = NewRepository(configuration)
	buildMenu(app)
}
