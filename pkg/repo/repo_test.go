package repo_test

import (
	"testing"

	"github.com/pixelfactoryio/git-get/pkg/repo"
	"github.com/stretchr/testify/require"
)

var tests []Test

type Test struct {
	in          string // Input string
	gitHost     string // parsed host
	projectPath string // parsed project path
	projectName string // parsed project name
	parsedURL   string // parsed project url
	isValid     bool   // is valid scheme
}

func newTest(in string, gitHost string, projectPath string, projectName string, parsedURL string, isValid bool) Test {

	if parsedURL == "" {
		parsedURL = in
	}

	return Test{
		in:          in,
		gitHost:     gitHost,
		projectPath: projectPath,
		projectName: projectName,
		parsedURL:   parsedURL,
		isValid:     isValid,
	}
}

func Test_New(t *testing.T) {
	t.Parallel()
	is := require.New(t)

	tests = []Test{
		newTest("http://githost.com/path/to/repo.git", "githost.com", "path/to", "repo", "", true),
		newTest("https://githost.com/path/to/repo.git", "githost.com", "path/to", "repo", "", true),
		newTest("ftp://githost.com/path/to/repo.git", "githost.com", "path/to", "repo", "", true),
		newTest("ftps://githost.com/path/to/repo.git", "githost.com", "path/to", "repo", "", true),
		newTest("ssh://git@githost.com/path/to/repo.git", "githost.com", "path/to", "repo", "", true),
		newTest("git+ssh://git@githost.com/path/to/repo.git", "githost.com", "path/to", "repo", "", true),
		newTest("git@githost.com:path/to/repo.git", "githost.com", "path/to", "repo", "ssh://git@githost.com/path/to/repo.git", true),
		newTest("unkown://path/to/repo.git", "githost.com", "path/to", "repo", "", false),
		newTest("httptt://path/to/repo.git", "githost.com", "path/to", "repo", "", false),
	}

	for _, tt := range tests {
		r, err := repo.New(tt.in)

		if !tt.isValid {
			is.Empty(r)
			is.Error(err)
		}

		if tt.isValid {
			is.NoError(err)
			is.NotEmpty(r)
			is.Equal(r.URL, tt.parsedURL)
			is.Equal(r.GitHost, tt.gitHost)
			is.Equal(r.ProjectPath, tt.projectPath)
			is.Equal(r.ProjectName, tt.projectName)
		}
	}
}
