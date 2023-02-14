// Package project is used to handle project
package project

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	giturls "github.com/whilp/git-urls"

	internalErrors "github.com/pixelfactoryio/git-get/internal/errors"
)

var gitURLRe = regexp.MustCompile(
	`(?m)(?:^git|^ssh|^https?|^git\+ssh|^git@[-\w.]+):(\/\/)?(.*?)(\.git)(\/?|\#[-\d\w._]+?)$`,
)

// Cloner interface
type Cloner interface {
	Clone(name *Project) error
}

// Project respresents a git project
type Project struct {
	URL      string
	GitHost  string
	Name     string
	Path     string
	FullPath string
}

// New creates a new project Repo using the given url
func New(rawRepoURL string, srcPath string) (*Project, error) {
	ok, err := isValidGitURL(rawRepoURL)
	if !ok {
		return nil, err
	}

	u, err := giturls.Parse(rawRepoURL)
	if err != nil {
		return nil, internalErrors.WrapErrorf(err, internalErrors.ErrorCodeInvalidArgument, "unable to parse URL")
	}

	basename := path.Base(u.Path)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))
	ppath := strings.TrimPrefix(path.Dir(u.Path), "/")
	fullPath := filepath.Join(srcPath, u.Host, ppath, name)

	r := &Project{
		URL:      u.String(),
		GitHost:  u.Host,
		Name:     name,
		Path:     ppath,
		FullPath: fullPath,
	}

	return r, nil
}

func isValidGitURL(rawRepoURL string) (bool, error) {
	match := gitURLRe.FindAllString(rawRepoURL, -1)
	if len(match) == 0 {
		return false, internalErrors.NewErrorf(
			internalErrors.ErrorCodeInvalidArgument,
			fmt.Sprintf("invalid URL scheme: %s", rawRepoURL),
		)
	}

	return true, nil
}
