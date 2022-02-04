package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
	ErrSignal     = errors.New("received signal")
)

type stepErr struct {
	step  string
	msg   string
	cause error
}

// Implement the error interface.
func (s *stepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

func (s *stepErr) Is(target error) bool {
	t, ok := target.(*stepErr)
	if !ok {
		return false
	}

	return t.step == s.step
}

// Unwrap allows errors.Is to check if underlying errors match the target.
func (s *stepErr) Unwrap() error {
	return s.cause
}
