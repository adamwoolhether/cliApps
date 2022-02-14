package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type executor interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("p", "", "Project directory")
	branch := flag.String("b", "", "Git branch")
	flag.Parse()
	
	if err := run(*proj, *branch, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj, branch string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}
	
	if branch == "" {
		s := bufio.NewScanner(os.Stdin)
		fmt.Println("Input Git repo's target branch")
		for s.Scan() {
			if err := s.Err(); err != nil {
				return err
			}
			if len(s.Text()) > 1 {
				branch = s.Text()
				break
			}
			fmt.Println("target branch can't be blank")
		}
	}
	
	pipeline := make([]executor, 6)
	
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj, []string{"build", ".", "errors"},
	)
	
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj, []string{"test", "-v"},
	)
	
	pipeline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		proj, []string{"-l", "."},
	)
	
	pipeline[3] = newExceptionStep(
		"ci-lint",
		"golangci-lint",
		"CI-Lint: SUCCESS",
		proj, []string{"run"},
	)
	
	pipeline[4] = newExceptionStep(
		"gocyclo",
		"gocyclo",
		"Gocyclo: SUCCESS",
		proj, []string{"-over", "10", "."})
	
	pipeline[5] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS",
		proj, []string{"push", "origin", branch},
		10*time.Second,
	)
	
	sig := make(chan os.Signal, 1)
	errCh := make(chan error)
	done := make(chan struct{})
	
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		for _, s := range pipeline {
			msg, err := s.execute()
			if err != nil {
				errCh <- err
				return
			}
			
			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				errCh <- err
				return
			}
		}
		close(done)
	}()
	
	for {
		select {
		case rec := <-sig:
			signal.Stop(sig)
			return fmt.Errorf("%s: exiting: %w", rec, ErrSignal)
		case err := <-errCh:
			return err
		case <-done:
			return nil
		}
	}
	
	return nil
}
