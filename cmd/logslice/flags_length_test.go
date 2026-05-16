package main

import (
	"testing"

	"github.com/user/logslice/internal/filter"
)

func TestBuildLengthFilter_NoOp(t *testing.T) {
	chain := filter.NewChain()
	if err := buildLengthFilter(chain, "", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// chain should remain empty
	if chain.Len() != 0 {
		t.Errorf("expected empty chain, got len %d", chain.Len())
	}
}

func TestBuildLengthFilter_InvalidOp(t *testing.T) {
	chain := filter.NewChain()
	if err := buildLengthFilter(chain, "neq", 5); err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestBuildLengthFilter_NegativeVal(t *testing.T) {
	chain := filter.NewChain()
	if err := buildLengthFilter(chain, "gt", -1); err == nil {
		t.Fatal("expected error for negative val")
	}
}

func TestBuildLengthFilter_Valid(t *testing.T) {
	chain := filter.NewChain()
	if err := buildLengthFilter(chain, "gte", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chain.Len() != 1 {
		t.Errorf("expected chain len 1, got %d", chain.Len())
	}
}
