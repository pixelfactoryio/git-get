package project

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	giturls "github.com/whilp/git-urls"

	"github.com/pixelfactoryio/git-get/internal"
)

var (
	gitUrlRe = regexp.MustCompile(`(?m)(?:^git|^ssh|^https?|^git\+ssh|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$`)
)

// ProjectRepository interface
type ProjectCloner interface {
	Clone(name *Repo) (string, error)
}

// Repo respresents a git repository
type Repo struct {
	URL             string
	GitHost         string
	ProjectName     string
	ProjectPath     string
	FullProjectPath string
}

// New creates a new project Repo using the given url
func New(rawRepoURL string, srcPath string) (*Repo, error) {

	ok, err := isValidGitURL(rawRepoURL)
	if !ok {
		return nil, err
	}

	u, err := giturls.Parse(rawRepoURL)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "unable to parse URL")
	}

	basename := path.Base(u.Path)
	projectName := strings.TrimSuffix(basename, filepath.Ext(basename))
	projectPath := strings.TrimPrefix(path.Dir(u.Path), "/")
	fullProjectPath := fmt.Sprintf("%s/%s/%s/%s", srcPath, u.Host, projectPath, projectName)

	r := &Repo{
		URL:             u.String(),
		GitHost:         u.Host,
		ProjectName:     projectName,
		ProjectPath:     projectPath,
		FullProjectPath: fullProjectPath,
	}

	return r, nil
}

func isValidGitURL(rawRepoURL string) (bool, error) {
	match := gitUrlRe.FindAllString(rawRepoURL, -1)
	if len(match) == 0 {
		return false, internal.NewErrorf(internal.ErrorCodeInvalidArgument, fmt.Sprintf("invalid URL scheme: %s", rawRepoURL))
	}

	return true, nil
}
