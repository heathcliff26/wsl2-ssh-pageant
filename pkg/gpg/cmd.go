package gpg

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	flagNameGPGConfigDir = "config-dir"
)

func NewCommand() (*cobra.Command, error) {
	userhome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	defaultGPGConfigDir := filepath.Join(userhome, "AppData", "Roaming", "gnupg")

	cmd := &cobra.Command{
		Use:   "gpg <gpg-agent-socket>",
		Short: "Connect GPG Agent socket on windows to unix socket via socat",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gpgConfigDir, err := cmd.Flags().GetString(flagNameGPGConfigDir)
			if err != nil {
				slog.Error("Unkown error parsing command", "err", err)
				os.Exit(1)
			}
			path := filepath.Join(gpgConfigDir, args[0])
			handleGPG(path)
		},
	}
	cmd.PersistentFlags().String(flagNameGPGConfigDir, defaultGPGConfigDir, "gpg config directory on windows")

	return cmd, nil
}
