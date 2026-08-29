package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type ExecCommand interface {
	Git(name string, args ...string) (string, error)
}

type execCommand struct {
	wd string
}

func NewExecutor(wd string) ExecCommand {
	return &execCommand{
		wd: wd,
	}
}

func (e *execCommand) Git(name string, args ...string) (string, error) {
	cmdArgs := []string{name}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Dir = e.wd

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("executing git command: %w", err)
	}
	trimmed := strings.TrimSpace(string(out))
	return trimmed, nil
}
