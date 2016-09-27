package version

import (
	"fmt"
)

var Version = "Development"
var Build = "No Build Number Given"
var Commit = "No Commit Hash Given"
var Branch = "No Branch Given"
var Repo = "No Repo Given"
var BuildTime = "Sometime in the Past"
var Tag = "No Tag Given"

func PrintVersion() string {
	return fmt.Sprintf("Version:\t%s\nTag:\t%s\nBuild:\t%sTime:\t%s\n\nCommit:\t%s\nRepo:\t%s\nBranch:\t%s\n",
		Version, Tag, Build, BuildTime, Commit, Repo, Branch,
	)
}
