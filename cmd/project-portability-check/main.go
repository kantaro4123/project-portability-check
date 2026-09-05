package main

import (
	"context"
	"os"

	"github.com/kantaro4123/project-portability-check/internal/cli"
)

func main() {
	os.Exit(cli.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}
