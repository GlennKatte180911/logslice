package main

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "text" {
		t.Errorf("default Format = %q, want \"text\"", cfg.Format)
	}
	if cfg.From != "" || cfg.To != "" || cfg.Pattern != "" || cfg.Levels != "" || cfg.Input != "" {
		t.Error("expected all optional flags to default to empty string")
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	args := []string{
		"-from", "2024-01-01T00:00:00Z",
		"-to", "2024-01-02T00:00:00Z",
		"-pattern", "timeout",
		"-levels", "ERROR,WARN",
		"-format", "json",
		"-input", "/var/log/app.log",
	}
	cfg, err := parseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.From != "2024-01-01T00:00:00Z" {
		t.Errorf("From = %q", cfg.From)
	}
	if cfg.To != "2024-01-02T00:00:00Z" {
		t.Errorf("To = %q", cfg.To)
	}
	if cfg.Pattern != "timeout" {
		t.Errorf("Pattern = %q", cfg.Pattern)
	}
	if cfg.Levels != "ERROR,WARN" {
		t.Errorf("Levels = %q", cfg.Levels)
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %q", cfg.Format)
	}
	if cfg.Input != "/var/log/app.log" {
		t.Errorf("Input = %q", cfg.Input)
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestBuildChain_NoFilters(t *testing.T) {
	cfg := Config{Format: "text"}
	chain, err := buildChain(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chain == nil {
		t.Error("expected non-nil chain")
	}
}

func TestBuildChain_InvalidPattern(t *testing.T) {
	cfg := Config{Pattern: "[invalid"}
	_, err := buildChain(cfg)
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestBuildChain_InvalidLevels(t *testing.T) {
	cfg := Config{Levels: "ERROR,,WARN"}
	_, err := buildChain(cfg)
	if err == nil {
		t.Fatal("expected error for blank level entry, got nil")
	}
}

func TestBuildChain_ValidPatternAndLevels(t *testing.T) {
	cfg := Config{
		Format:  "text",
		Pattern: "timeout",
		Levels:  "ERROR,WARN,INFO",
	}
	chain, err := buildChain(cfg)
	if err != nil {
		t.Fatalf("unexpected error building chain with pattern and levels: %v", err)
	}
	if chain == nil {
		t.Error("expected non-nil chain")
	}
}
