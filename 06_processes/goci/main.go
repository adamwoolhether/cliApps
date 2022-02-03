package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	// Use go build to validate compilation.
	// We add a lib in stdlib to build to prevent building a binary to clean up.
	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)
	cmd.Dir = proj
	if err := cmd.Run(); err != nil {
		return &stepErr{step: "go build", msg: "go build failed", cause: err}
	}

	_, err := fmt.Println(out, "Go Build: SUCCESS")

	return err
}
