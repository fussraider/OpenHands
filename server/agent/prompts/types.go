package prompts

type RepositoryInfo struct {
	RepoName      string
	RepoDirectory string
	BranchName    string
}

type RuntimeInfo struct {
	WorkingDir                  string
	AvailableHosts              map[string]int
	AdditionalAgentInstructions string
	CustomSecretsDescriptions   map[string]string
	Date                        string
}

type ConversationInstructions struct {
	Content string
}

type SystemPromptContext struct {
	CLIMode bool
}

type AdditionalInfoContext struct {
	RepositoryInfo           *RepositoryInfo
	RepositoryInstructions   string
	RuntimeInfo              *RuntimeInfo
	ConversationInstructions *ConversationInstructions
}
