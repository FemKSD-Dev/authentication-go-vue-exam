package migrations

import (
	"testing"
)

func TestFiles_Embedded(t *testing.T) {
	// embed.FS is never nil; verify we can read from it
	entries, err := Files.ReadDir(".")
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(entries) == 0 {
		t.Log("no sql files in migrations (may be empty)")
	}
}
