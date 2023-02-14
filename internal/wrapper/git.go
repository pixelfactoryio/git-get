// Package wrapper is used to wrap git commands
package wrapper

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ldez/go-git-cmd-wrapper/v2/clone"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/types"

	internalErrors "github.com/pixelfactoryio/git-get/internal/errors"
	"github.com/pixelfactoryio/git-get/internal/project"
)

type gitCmd struct {
	cmdExecutor types.Executor
}

// Option is an option for New git wrapper.
type Option func(*gitCmd)

// WithCmdExecutor set git-wrapper executor
func WithCmdExecutor(e types.Executor) Option {
	return func(g *gitCmd) {
		g.cmdExecutor = e
	}
}

func cmdExecutor(ctx context.Context, name string, _ bool, args ...string) (string, error) {
	gitBin, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, gitBin, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return "", err
}

// NewGitWrapper returns wrapper around git command
func NewGitWrapper(opts ...Option) project.Cloner {
	g := &gitCmd{
		cmdExecutor: cmdExecutor,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func (g *gitCmd) Clone(r *project.Project) error {
	empty, err := isEmptyDirOrNotExist(r.FullPath)
	if err != nil || !empty {
		return internalErrors.NewErrorf(
			internalErrors.ErrorCodeInvalidArgument, "destination path %s is not empty", r.FullPath,
		)
	}

	err = os.MkdirAll(r.FullPath, 0o755) //nolint:gosec
	if err != nil {
		return err
	}

	_, err = git.Clone(
		clone.Repository(r.URL), clone.Directory(r.FullPath), clone.Progress, git.CmdExecutor(g.cmdExecutor),
	)
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return internalErrors.NewErrorf(internalErrors.ErrorCode(ee.ExitCode()), "git client returned error")
		}

		return internalErrors.WrapErrorf(err, internalErrors.ErrorCodeUnknown, "git client returned error")
	}

	return nil
}

func isEmptyDirOrNotExist(name string) (found bool, err error) {
	filePath := filepath.Clean(name)
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	if err = f.Close(); err != nil {
		return true, err
	}

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err // Either not empty or error, suits both cases
}
