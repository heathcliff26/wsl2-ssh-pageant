package version

import (
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

const Name = "wsl2-ssh-pageant"

// Create a new version command with the given app name
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information and exit",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(VersionInfoString())
		},
	}
	// Override to prevent parent function from running
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}

	return cmd
}

// Return the version string
func Version() string {
	buildinfo, _ := debug.ReadBuildInfo()
	return buildinfo.Main.Version
}

// Return a formated string containing the version, git commit and go version the app was compiled with.
func VersionInfoString() string {
	var commit string
	buildinfo, _ := debug.ReadBuildInfo()
	for _, item := range buildinfo.Settings {
		if item.Key == "vcs.revision" {
			commit = item.Value
			break
		}
	}
	if len(commit) > 7 {
		commit = commit[:7]
	} else if commit == "" {
		commit = "Unknown"
	}

	result := Name + ":\n"
	result += "    Version: " + buildinfo.Main.Version + "\n"
	result += "    Commit:  " + commit + "\n"
	result += "    Go:      " + runtime.Version() + "\n"

	return result
}
