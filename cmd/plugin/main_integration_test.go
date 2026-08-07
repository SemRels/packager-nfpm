package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestIntegrationDryRunBuildPlanAndInvalidVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, packagerEnv(map[string]string{
		"SEMREL_VERSION":          "1.2.3",
		"SEMREL_PLUGIN_CONFIG":    "nfpm.yaml",
		"SEMREL_PLUGIN_TARGET":    "dist",
		"SEMREL_PLUGIN_PACKAGERS": "deb,rpm",
		"SEMREL_DRY_RUN":          "true",
	}))
	if code != 0 {
		t.Fatalf("dry-run code = %d, stderr = %q", code, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, "nfpm package --config nfpm.yaml --target dist --packager deb") ||
		!strings.Contains(got, "nfpm package --config nfpm.yaml --target dist --packager rpm") {
		t.Errorf("dry-run plan = %q", got)
	}

	stdout.Reset()
	stderr.Reset()
	code = run(&stdout, &stderr, packagerEnv(nil))
	if code != 1 || !strings.Contains(stderr.String(), "SEMREL_VERSION is required") {
		t.Fatalf("invalid configuration: code = %d, stderr = %q", code, stderr.String())
	}
}

func packagerEnv(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
