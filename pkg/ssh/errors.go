package ssh

import "fmt"

type ErrMessageLength struct {
	length int
}

func NewErrMessageLength(l int) error {
	return &ErrMessageLength{
		length: l,
	}
}

func (e *ErrMessageLength) Error() string {
	return fmt.Sprintf("Message too long, max length is %d, current message length is %d", agentMaxMessageLength, e.length)
}

type ErrNoPageant struct{}

func NewErrNoPageant() error {
	return ErrNoPageant{}
}

func (e ErrNoPageant) Error() string {
	return "Could not find Pageant window"
}

type ErrWMCopydata struct{}

func NewErrWMCopydata() error {
	return ErrWMCopydata{}
}

func (e ErrWMCopydata) Error() string {
	return "WM_COPYDATA failed"
}

type ErrGPGAgentLaunch struct {
	err error
}

func NewErrGPGAgentLaunch(err error) error {
	return &ErrGPGAgentLaunch{
		err: err,
	}
}

func (e *ErrGPGAgentLaunch) Error() string {
	return fmt.Sprintf("Failed to launch gpg-connect-agent: %v", e.err)
}
