package wrapper

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/ldez/go-git-cmd-wrapper/v2/clone"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"

	"github.com/pixelfactoryio/git-get/internal"
	"github.com/pixelfactoryio/git-get/internal/project"
)

type gitCmd struct{}

func NewGitWrapper() project.ProjectCloner {
	return &gitCmd{}
}

func (g *gitCmd) Clone(r *project.Repo) (string, error) {
	empty, err := isEmptyDirOrNotExist(r.FullProjectPath)
	if err != nil || !empty {
		return "", internal.NewErrorf(internal.ErrorCodeInvalidArgument, "destination path %s is not empty", r.FullProjectPath)
	}

	err = os.MkdirAll(r.FullProjectPath, 0755)
	if err != nil {
		return "", err
	}

	out, err := git.Clone(clone.Repository(r.URL), clone.Directory(r.FullProjectPath), clone.Progress)
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return "", internal.NewErrorf(uint(ee.ExitCode()), "git client returned error:\n%v", out)
		}

		return "", internal.WrapErrorf(err, internal.ErrorCodeUnknown, "git client returned error:\n%v", out)
	}

	return out, nil
}

func isEmptyDirOrNotExist(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err // Either not empty or error, suits both cases
}
