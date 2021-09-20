package repo

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	giturls "github.com/whilp/git-urls"

	"github.com/pixelfactoryio/git-get/internal"
)

var transports = []string{
	"ssh://",
	"git://",
	"git+ssh://",
	"http://",
	"https://",
	"ftp://",
	"ftps://",
	"rsync://",
	"file://",
	"git@",
}

// Repo respresents a git repository
type Repo struct {
	GitHost     string
	ProjectName string
	ProjectPath string
	URL         string
}

// New creates a new git Repo using the given url
func New(rawRepoURL string) (*Repo, error) {

	ok, err := isValidScheme(rawRepoURL)
	if !ok {
		return nil, err
	}

	u, err := giturls.Parse(rawRepoURL)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "unable to parse URL")
	}

	basename := path.Base(u.Path)
	r := &Repo{
		URL:         u.String(),
		GitHost:     u.Host,
		ProjectName: strings.TrimSuffix(basename, filepath.Ext(basename)),
		ProjectPath: strings.TrimPrefix(path.Dir(u.Path), "/"),
	}

	return r, nil
}

func isValidScheme(rawRepoURL string) (bool, error) {
	for _, t := range transports {
		re := fmt.Sprintf("(?m)^%s", regexp.QuoteMeta(t))
		match, err := regexp.MatchString(re, rawRepoURL)
		if err != nil {
			return false, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, fmt.Sprintf("invalid URL scheme: %s", rawRepoURL))
		}

		if match {
			return true, nil
		}
	}
	return false, internal.NewErrorf(internal.ErrorCodeInvalidArgument, fmt.Sprintf("invalid URL scheme: %s", rawRepoURL))
}
