package cmd

import (
	"fmt"
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
	rootCmd := &cobra.Command{
		Use:           "git-get <repo>",
		Short:         "git-get",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE:       preStart,
		RunE:          start,
	}

	rootCmd.PersistentFlags().String("projects-path", "", "source directory path")
	err := viper.BindPFlag("projects-path", rootCmd.PersistentFlags().Lookup("projects-path"))
	if err != nil {
		return nil, err
	}

	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error, fatal, panic)")
	err = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		return nil, err
	}

	rootCmd.PersistentFlags().Bool("debug", false, "run in debug mode")
	err = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		return nil, err
	}

	return rootCmd, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd, err := NewRootCmd()
	if err != nil {
		return err
	}

	cobra.OnInitialize(initConfig)
	return rootCmd.Execute()
}

func preStart(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return internal.NewErrorf(internal.ErrorCodeInvalidArgument, "the <repo> argument is required")
	}

	if len(args) > 1 {
		return internal.NewErrorf(internal.ErrorCodeInvalidArgument, "the <repo> argument is required")
	}
	return nil
}

func start(c *cobra.Command, args []string) error {
	srcPath := viper.GetString("projects-path")
	if len(srcPath) == 0 {
		return internal.NewErrorf(internal.ErrorCodeInvalidArgument, "Please set GIT_GET_PROJECTS_PATH or use `--projects-path` flag")
	}

	r, err := project.New(args[0], srcPath)
	if err != nil {
		return err
	}

	git := wrapper.NewGitWrapper()
	out, err := git.Clone(r)
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}
