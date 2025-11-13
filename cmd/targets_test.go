package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTargetSpecWithDisplay(t *testing.T) {
	cfg, err := parseTargetSpec("api=example.com:443", "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.display != "api" {
		t.Fatalf("expected display api, got %s", cfg.display)
	}
	if cfg.address != "example.com:443" {
		t.Fatalf("expected address example.com:443, got %s", cfg.address)
	}
}

func TestParseTargetSpecDefaultDisplay(t *testing.T) {
	cfg, err := parseTargetSpec("example.com:80", "udp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.display != "udp://example.com:80" {
		t.Fatalf("unexpected display %s", cfg.display)
	}
}

func TestParseTargetSpecErrors(t *testing.T) {
	if _, err := parseTargetSpec("missing-port", "tcp"); err == nil {
		t.Fatalf("expected error for missing port")
	}
	if _, err := parseTargetSpec("name=", "tcp"); err == nil {
		t.Fatalf("expected error for empty host")
	}
}

func TestMultiTargetFlag(t *testing.T) {
	var flag multiTargetFlag
	if err := flag.Set("one=host:1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := flag.Set("two:2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	values := flag.Values()
	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}
	if values[0] != "one=host:1" || values[1] != "two:2" {
		t.Fatalf("unexpected values: %#v", values)
	}
}

func TestReadTargetsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "targets.txt")
	content := "# comment\nexample.com:80\n\nname=example.org:443"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	entries, err := readTargetsFile(path)
	if err != nil {
		t.Fatalf("readTargetsFile err: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0] != "example.com:80" || entries[1] != "name=example.org:443" {
		t.Fatalf("unexpected entries: %#v", entries)
	}
}
