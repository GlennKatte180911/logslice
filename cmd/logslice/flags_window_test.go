package main

import (
	"testing"
)

func TestBuildWindowFilter_NoAnchor(t *testing.T) {
	cfg := &config{}
	f, err := buildWindowFilter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Fatal("expected nil filter when no anchor set")
	}
}

func TestBuildWindowFilter_InvalidAnchor(t *testing.T) {
	cfg := &config{
		windowAnchor: "not-a-time",
		windowBefore: "5m",
		windowAfter:  "5m",
	}
	_, err := buildWindowFilter(cfg)
	if err == nil {
		t.Fatal("expected error for invalid anchor")
	}
}

func TestBuildWindowFilter_InvalidBefore(t *testing.T) {
	cfg := &config{
		windowAnchor: "2024-06-01T12:00:00Z",
		windowBefore: "bad",
		windowAfter:  "5m",
	}
	_, err := buildWindowFilter(cfg)
	if err == nil {
		t.Fatal("expected error for invalid before duration")
	}
}

func TestBuildWindowFilter_InvalidAfter(t *testing.T) {
	cfg := &config{
		windowAnchor: "2024-06-01T12:00:00Z",
		windowBefore: "5m",
		windowAfter:  "bad",
	}
	_, err := buildWindowFilter(cfg)
	if err == nil {
		t.Fatal("expected error for invalid after duration")
	}
}

func TestBuildWindowFilter_Valid(t *testing.T) {
	cfg := &config{
		windowAnchor: "2024-06-01T12:00:00Z",
		windowBefore: "10m",
		windowAfter:  "5m",
	}
	f, err := buildWindowFilter(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}
