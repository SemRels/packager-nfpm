//go:build integration

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationBuildsInspectableDebianPackage(t *testing.T) {
	if _, err := exec.LookPath("nfpm"); err != nil {
		t.Fatalf("nfpm is required for the package integration test: %v", err)
	}
	if _, err := exec.LookPath("dpkg-deb"); err != nil {
		t.Fatalf("dpkg-deb is required for the package integration test: %v", err)
	}

	dir := t.TempDir()
	payload := filepath.Join(dir, "payload.txt")
	if err := os.WriteFile(payload, []byte("semrel package integration payload\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(dir, "nfpm.yaml")
	contents := strings.ReplaceAll(payload, "\\", "/")
	configBody := strings.Join([]string{
		"name: semrel-integration",
		"arch: amd64",
		"platform: linux",
		"version: 1.2.3",
		"section: default",
		"priority: extra",
		"maintainer: Semrel <release@example.test>",
		"description: Semrel package integration test",
		"contents:",
		"  - src: " + contents,
		"    dst: /usr/share/semrel/payload.txt",
	}, "\n") + "\n"
	if err := os.WriteFile(config, []byte(configBody), 0o600); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(dir, "semrel-integration.deb")
	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, packagerEnv(map[string]string{
		"SEMREL_VERSION":          "1.2.3",
		"SEMREL_PLUGIN_CONFIG":    config,
		"SEMREL_PLUGIN_TARGET":    target,
		"SEMREL_PLUGIN_PACKAGERS": "deb",
	}))
	if code != 0 {
		t.Fatalf("package code = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
	}

	if _, err := os.Stat(target); err != nil {
		t.Fatalf("built package %q: %v", target, err)
	}
	output, err := exec.Command("dpkg-deb", "--field", target, "Package").CombinedOutput()
	if err != nil {
		t.Fatalf("inspect package name: %v: %s", err, output)
	}
	if got := strings.TrimSpace(string(output)); got != "semrel-integration" {
		t.Fatalf("package name = %q, want %q", got, "semrel-integration")
	}

	output, err = exec.Command("dpkg-deb", "--field", target, "Version").CombinedOutput()
	if err != nil {
		t.Fatalf("inspect package version: %v: %s", err, output)
	}
	if got := strings.TrimSpace(string(output)); got != "1.2.3" {
		t.Fatalf("package version = %q, want %q", got, "1.2.3")
	}

	extractDir := filepath.Join(dir, "extracted")
	output, err = exec.Command("dpkg-deb", "--extract", target, extractDir).CombinedOutput()
	if err != nil {
		t.Fatalf("extract package: %v: %s", err, output)
	}
	extracted, err := os.ReadFile(filepath.Join(extractDir, "usr", "share", "semrel", "payload.txt"))
	if err != nil {
		t.Fatalf("read packaged payload: %v", err)
	}
	if string(extracted) != "semrel package integration payload\n" {
		t.Fatalf("packaged payload = %q", extracted)
	}
}
