package ssh

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

const (
	flagNamePipe = "pipe"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Connect Pageant to unix socket via socat",
		Run: func(cmd *cobra.Command, args []string) {
			pipe, err := cmd.Flags().GetString(flagNamePipe)
			if err != nil {
				slog.Error("Unkown error parsing command", "err", err)
				os.Exit(1)
			}
			if pipe == "" {
				handleSSH()
			} else {
				handlePipedSSH(pipe)
			}
		},
	}
	cmd.PersistentFlags().String(flagNamePipe, "", "Use a pipe for communicating with pageant instead of the default shared memory API")

	return cmd
}
