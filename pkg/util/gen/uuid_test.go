package gen

import "testing"

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	if uuid == "" {
		t.Fatalf("Generated UUID should not be empty")
	}
}