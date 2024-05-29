package entity

type Configuration struct {
	Repositories              []string `json:"repositories"`
	TerminalCommand           string   `json:"terminal-command"`
	LazygitCommand            string   `json:"lazygit-command"`
	MaxCharRepoRepresentation int      `json:"max-char-repo-representation"`
}
