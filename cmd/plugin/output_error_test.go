package main

import (
	"errors"
	"testing"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestRunDryRunSucceedsWhenOutputFails(t *testing.T) {
	code := run(failingWriter{}, failingWriter{}, packagerEnv(map[string]string{
		"SEMREL_VERSION": "1.2.3",
		"SEMREL_DRY_RUN": "true",
	}))
	if code != 0 {
		t.Fatalf("run() code = %d, want 0", code)
	}
}
