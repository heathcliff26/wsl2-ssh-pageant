package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/heathcliff26/wsl2-ssh-pageant/pkg/gpg"
	"github.com/heathcliff26/wsl2-ssh-pageant/pkg/ssh"
	"github.com/heathcliff26/wsl2-ssh-pageant/pkg/version"
	"github.com/spf13/cobra"
)

const (
	FlagNameLogfile  = "log-file"
	FlagNameLogLevel = "log-level"

	DefaultLogFile  = "wsl2-ssh-pageant.log"
	DefaultLogLevel = "info"
)

var (
	logFile  *os.File
	logLevel = &slog.LevelVar{}
)

func NewRootCommand() *cobra.Command {
	cobra.AddTemplateFunc(
		"ProgramName", func() string {
			return version.Name
		},
	)

	rootCmd := &cobra.Command{
		Use:              version.Name,
		Short:            version.Name + " allows using pageant inside WSL2 Distros as ssh-agent.",
		PersistentPreRun: preRun,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			os.Args = append([]string{os.Args[0], ssh.CommandName}, os.Args[1:]...)
			cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}
			err := cmd.Execute()
			if err != nil {
				exitError(cmd, err)
			}
		},
	}

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.PersistentFlags().String(FlagNameLogfile, DefaultLogFile, "Path to logfile")
	rootCmd.PersistentFlags().String(FlagNameLogLevel, DefaultLogLevel, "Log level")

	gpgCommand, err := gpg.NewCommand()
	if err != nil {
		exitError(rootCmd, err)
	}

	rootCmd.AddCommand(
		ssh.NewCommand(),
		gpgCommand,
		version.NewCommand(),
	)

	return rootCmd
}

func Execute() {
	defer postRun()

	cmd := NewRootCommand()
	err := cmd.Execute()
	if err != nil {
		exitError(cmd, err)
	}
}

// Initialization function run before every command
func preRun(cmd *cobra.Command, args []string) {
	logFilePath, err := cmd.Flags().GetString(FlagNameLogfile)
	if err != nil {
		exitError(cmd, err)
	}
	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		exitError(cmd, err)
	}
	logFile = f

	opts := slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(f, &opts))
	slog.SetDefault(logger)

	level, err := cmd.Flags().GetString(FlagNameLogLevel)
	if err != nil {
		exitError(cmd, err)
	}
	err = setLogLevel(level)
	if err != nil {
		exitError(cmd, err)
	}
}

// Cleanup log resources after program ends
func postRun() {
	if logFile != nil {
		err := logFile.Sync()
		if err != nil {
			fmt.Printf("Failed to flush logfile: %v", err)
		}
		logFile.Close()
		logFile = nil
	}
}

// Parse a given string and set the resulting log level
func setLogLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		return NewErrUnknownLogLevel(level)
	}
	return nil
}

// Print the error information on stderr and exit with code 1
func exitError(cmd *cobra.Command, err error) {
	fmt.Fprintln(cmd.Root().ErrOrStderr(), err.Error())
	os.Exit(1)
}
