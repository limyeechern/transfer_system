package util

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestGenerateTxID(t *testing.T) {
	first := GenerateTxID()
	second := GenerateTxID()

	if first == "" {
		t.Fatalf("expected tx id")
	}
	if _, err := ulid.Parse(first); err != nil {
		t.Fatalf("tx id is not a valid ULID: %v", err)
	}
	if first == second {
		t.Fatalf("expected unique tx ids, got %q twice", first)
	}
}
