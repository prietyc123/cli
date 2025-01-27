package mock

import (
	"bytes"
	"io"
	"net/url"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// FakeRemoteExecutor implements exec.RemoteExecutor interface, part of `kubectl exec` subcommand,
// capturing Execute method input.
type FakeRemoteExecutor struct {
	command []string     // extracted from query parameter ("command")
	stdin   bytes.Buffer // standard input informed
	err     error        // stubbed error
}

// Command returns the command informed to Execute.
func (f *FakeRemoteExecutor) Command() []string {
	return f.command
}

// Stdin returns as string the stdin informed to Execute.
func (f *FakeRemoteExecutor) Stdin() string {
	return f.stdin.String()
}

// Execute handles the actual http request against Kubernetes API, and here greatly simplified to
// only return a stubbed error, and extract elements from the request.
func (f *FakeRemoteExecutor) Execute(
	_ string,
	reqURL *url.URL,
	_ *rest.Config,
	stdin io.Reader,
	stdout, stderr io.Writer,
	_ bool,
	_ remotecommand.TerminalSizeQueue,
) error {
	values, exists := reqURL.Query()["command"]
	if exists {
		f.command = values
	}

	if stdin != nil {
		if _, err := io.Copy(&f.stdin, stdin); err != nil {
			return err
		}
	}
	return f.err
}

// NewFakeRemoteExecutor instantiate a FakeRemoteExecutor with stubbed error.
func NewFakeRemoteExecutor(err error) *FakeRemoteExecutor {
	return &FakeRemoteExecutor{
		command: []string{},
		err:     err,
	}
}
