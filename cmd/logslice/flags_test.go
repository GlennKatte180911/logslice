package main

import (
	"strings"
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.from != "" {
		t.Errorf("expected empty from, got %q", cfg.from)
	}
	if cfg.to != "" {
		t.Errorf("expected empty to, got %q", cfg.to)
	}
	if cfg.pattern != "" {
		t.Errorf("expected empty pattern, got %q", cfg.pattern)
	}
	if cfg.input == nil {
		t.Error("expected non-nil input (stdin)")
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	cfg, err := parseFlags([]string{
		"-from", "2024-01-01T00:00:00Z",
		"-to", "2024-01-02T00:00:00Z",
		"-pattern", "ERROR",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.from != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected from: %q", cfg.from)
	}
	if cfg.to != "2024-01-02T00:00:00Z" {
		t.Errorf("unexpected to: %q", cfg.to)
	}
	if cfg.pattern != "ERROR" {
		t.Errorf("unexpected pattern: %q", cfg.pattern)
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Error("expected error for unknown flag, got nil")
	}
}

func TestBuildChain_NoFilters(t *testing.T) {
	cfg := &config{input: strings.NewReader("")}
	chain, err := buildChain(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chain == nil {
		t.Error("expected non-nil chain")
	}
}

func TestBuildChain_InvalidPattern(t *testing.T) {
	cfg := &config{pattern: "[invalid", input: strings.NewReader("")}
	_, err := buildChain(cfg)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestBuildChain_InvalidTimeRange(t *testing.T) {
	cfg := &config{from: "not-a-time", input: strings.NewReader("")}
	_, err := buildChain(cfg)
	if err == nil {
		t.Error("expected error for invalid time range")
	}
}
