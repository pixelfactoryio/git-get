package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.pixelfactory.io/pkg/version"

	"github.com/pixelfactoryio/git-get/internal"
	"github.com/pixelfactoryio/git-get/internal/project"
	"github.com/pixelfactoryio/git-get/internal/wrapper"
)

func initConfig() {
	viper.Set("revision", version.REVISION)
	viper.SetEnvPrefix("GIT_GET")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

// NewRootCmd create new rootCmd
func NewRootCmd() (*cobra.Command, error) {
	c := &cobra.Command{
		Use:           "git-get <repo>",
		Short:         "git-get",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE:       preStart,
		RunE:          start,
	}

	c.PersistentFlags().String("projects-path", "", "projects directory path")
	err := viper.BindPFlag("projects-path", c.PersistentFlags().Lookup("projects-path"))
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Execute adds all child commands to the root command
// and sets flags appropriately.
func Execute() error {
	c, err := NewRootCmd()
	if err != nil {
		return err
	}

	cobra.OnInitialize(initConfig)
	return c.Execute()
}

func preStart(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return internal.NewErrorf(internal.ErrorMissingArgument, "the <repo> argument is required")
	}

	if len(args) > 1 {
		return internal.NewErrorf(internal.ErrorMissingArgument, "the <repo> argument is required")
	}
	return nil
}

func start(c *cobra.Command, args []string) error {
	path := viper.GetString("projects-path")
	if len(path) == 0 {
		return internal.NewErrorf(
			internal.ErrorMissingArgument,
			"please set GIT_GET_PROJECTS_PATH or use `--projects-path` flag",
		)
	}

	p, err := project.New(args[0], path)
	if err != nil {
		return err
	}

	return wrapper.NewGitWrapper().Clone(p)
}
