package git

import (
	"fmt"
	"os"
	"strings"
)

type Repo struct {
	workingDir string
	executor   ExecCommand
}

func (r *Repo) Git(name string, args ...string) (string, error) {
	return r.executor.Git(name, args...)
}

func (r *Repo) WorkingDir() string {
	return r.workingDir
}

func LoadRepo() (*Repo, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	executor := NewExecutor(wd)
	output, err := executor.Git("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}
	if strings.TrimSpace(output) != "true" {
		return nil, fmt.Errorf("not a git repository: unexpected output %q", output)
	}

	return &Repo{workingDir: wd, executor: executor}, nil
}
