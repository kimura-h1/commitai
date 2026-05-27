// Package git provides helpers to interact with the local git repository.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// StagedDiff returns the output of `git diff --staged`.
// An error is returned if git is not available or the command fails.
func StagedDiff() (string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("git", "diff", "--staged")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git diff --staged: %s", msg)
	}

	return stdout.String(), nil
}

// Commit runs `git commit -m <message>` in the current directory.
func Commit(message string) error {
	var stderr bytes.Buffer

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("git commit: %s", msg)
	}

	return nil
}
