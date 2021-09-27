package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/pixelfactoryio/git-get/cmd"
	"github.com/pixelfactoryio/git-get/internal"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
		var ee *internal.Error
		if errors.As(err, &ee) {
			os.Exit(int(ee.Code()))
		}
		os.Exit(1)
	}
}
