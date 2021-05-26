package main

import (
	"fmt"
	"github.com/chaudhryfaisal/k8s-webhook-pull-policy/cmd"
	"os"
)

func main() {
	err := cmd.RunApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running app: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
