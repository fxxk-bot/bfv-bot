package config

var (
	GitCommit string
	GitBranch string
	BuildTime string
)

func GetVersion() string {
	return "Git Branch: " + GitBranch + "\n" +
		"Git Commit: " + GitCommit + "\n" +
		"Build Time: " + BuildTime
}
