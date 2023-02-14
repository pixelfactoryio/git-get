package wrapper_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pixelfactoryio/git-get/internal/project"
	"github.com/pixelfactoryio/git-get/internal/wrapper"
)

func cmdExecutorMock(_ context.Context, name string, _ bool, args ...string) (string, error) {
	return fmt.Sprintln(name, strings.Join(args, " ")), nil
}

func cmdExecutorMockError(ctx context.Context, name string, _ bool, args ...string) (string, error) {
	output, err := exec.CommandContext(ctx, "./fakegit.sh", args...).CombinedOutput()
	return string(output), err
}

func Test_New(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	g := wrapper.NewGitWrapper(wrapper.WithCmdExecutor(cmdExecutorMock))
	is.NotEmpty(g)
	is.Implements((*project.Cloner)(nil), g)
}

func Test_Clone(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	g := wrapper.NewGitWrapper(wrapper.WithCmdExecutor(cmdExecutorMock))
	p, err := project.New("http://githost.com/path/to/repo.git", t.TempDir())
	if err != nil {
		t.Error(err)
	}

	err = g.Clone(p)
	is.NoError(err)
}

func Test_Clone_DirExists(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	srcPath := t.TempDir()
	fullProjectPath := fmt.Sprintf("%s/githost.com/path/to/repo", srcPath)
	err := os.MkdirAll(fullProjectPath, 0o755) //nolint:gosec
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/test.txt", fullProjectPath), []byte("test file"), 0o755) //nolint:gosec
	if err != nil {
		t.Error(err)
	}

	g := wrapper.NewGitWrapper(wrapper.WithCmdExecutor(cmdExecutorMock))
	p, err := project.New("http://githost.com/path/to/repo.git", srcPath)
	if err != nil {
		t.Error(err)
	}

	err = g.Clone(p)
	is.Error(err)
}

func Test_Clone_GitError(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	g := wrapper.NewGitWrapper(wrapper.WithCmdExecutor(cmdExecutorMockError))
	p, err := project.New("http://githost.com/path/to/repo.git", t.TempDir())
	if err != nil {
		t.Error(err)
	}

	err = g.Clone(p)
	is.Error(err)
}
