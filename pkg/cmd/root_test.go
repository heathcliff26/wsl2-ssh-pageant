package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/heathcliff26/wsl2-ssh-pageant/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()

	assert.Equal(t, version.Name, cmd.Use)
}

func TestSetupAndCleanup(t *testing.T) {
	cmd := NewRootCommand()

	_ = cmd.ParseFlags([]string{})

	preRun(cmd, []string{})

	assert := assert.New(t)

	assert.NotNil(logFile)
	assert.Equal(DefaultLogLevel, strings.ToLower(logLevel.Level().String()))

	postRun()

	assert.Nil(logFile)

	err := os.Remove(DefaultLogFile)
	if err != nil {
		t.Fatalf("Failed to remove logfile: %v", err)
	}
}
