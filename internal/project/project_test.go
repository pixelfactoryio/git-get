package project_test

import (
	"testing"

	"github.com/pixelfactoryio/git-get/internal/project"
	"github.com/stretchr/testify/require"
)

var tests []Test

type Test struct {
	in          string // Input string
	gitHost     string // parsed host
	projectPath string // parsed project path
	projectName string // parsed project name
	parsedURL   string // parsed project url (if empty will default to in value)
	isValid     bool   // is valid scheme
}

func Test_New(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	tests = []Test{
		{
			"http://githost.com/path/to/repo.git",
			"githost.com",
			"path/to",
			"repo",
			"http://githost.com/path/to/repo.git",
			true,
		},
		{
			"https://githost.com/path/to/repo.git",
			"githost.com",
			"path/to",
			"repo",
			"https://githost.com/path/to/repo.git",
			true,
		},
		{
			"ssh://git@githost.com/path/to/repo.git",
			"githost.com",
			"path/to",
			"repo",
			"ssh://git@githost.com/path/to/repo.git",
			true,
		},
		{
			"git+ssh://git@githost.com/path/to/repo.git",
			"githost.com",
			"path/to",
			"repo",
			"git+ssh://git@githost.com/path/to/repo.git",
			true,
		},
		{
			"git@githost.com:path/to/repo.git",
			"githost.com",
			"path/to",
			"repo",
			"ssh://git@githost.com/path/to/repo.git",
			true,
		},
		{
			"unknown://githost.com/path/to/repo.git",
			"githost.com",
			"path/to",
			"repo",
			"unknown://githost.com/path/to/repo.git",
			false,
		},
		{
			"httptt://githost.com/path/to/repo.git",
			"githost.com",
			"path/to",
			"repo",
			"httptt://githost.com/path/to/repo.git",
			false,
		},
	}

	projectsPath := t.TempDir()
	for _, tt := range tests {
		p, err := project.New(tt.in, projectsPath)

		if !tt.isValid {
			is.Empty(p)
			is.Error(err)
		}

		if tt.isValid {
			is.NoError(err)
			is.NotEmpty(p)
			is.Equal(p.URL, tt.parsedURL)
			is.Equal(p.GitHost, tt.gitHost)
			is.Equal(p.Path, tt.projectPath)
			is.Equal(p.Name, tt.projectName)
		}
	}
}
