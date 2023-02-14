// Package main is the main entrypoint of this program
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/pixelfactoryio/git-get/cmd"
	internalErrors "github.com/pixelfactoryio/git-get/internal/errors"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
		var ee *internalErrors.Error
		if errors.As(err, &ee) {
			os.Exit(int(ee.Code()))
		}
		os.Exit(1)
	}
}
