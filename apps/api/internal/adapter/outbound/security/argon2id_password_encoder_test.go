package security

import (
	"errors"
	"strings"
	"testing"
)

func TestArgon2idPasswordEncoder_EncodeAndMatch(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	hash, err := enc.Encode("super-secret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if err := enc.Matches(hash, "super-secret"); err != nil {
		t.Fatalf("expected password to match, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_Mismatch(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	hash, err := enc.Encode("super-secret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := enc.Matches(hash, "wrong-password"); err == nil {
		t.Fatal("expected mismatch error, got nil")
	}
}

func TestArgon2idPasswordEncoder_EncodeSamePassword_ShouldProduceDifferentHashes(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	hash1, err := enc.Encode("same-password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	hash2, err := enc.Encode("same-password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("expected hashes to differ due to different salts")
	}
}

func TestArgon2idPasswordEncoder_InvalidHash(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	err := enc.Matches("not-a-valid-hash", "secret")
	if err == nil {
		t.Fatal("expected error for invalid hash")
	}
	if err != ErrorInvalidHash {
		t.Fatalf("expected ErrorInvalidHash, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_InvalidHash_WrongPartsCount(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	err := enc.Matches("$argon2id$v=19", "secret")
	if err == nil {
		t.Fatal("expected error for invalid hash")
	}
	if err != ErrorInvalidHash {
		t.Fatalf("expected ErrorInvalidHash, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_InvalidHash_WrongAlgorithm(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	err := enc.Matches("$argon2i$v=19$m=19456,t=2,p=1$AAAAAAAAAAAAAAAAAAAAAA==$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", "secret")
	if err == nil {
		t.Fatal("expected error for wrong algorithm")
	}
	if err != ErrorInvalidHash {
		t.Fatalf("expected ErrorInvalidHash, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_InvalidHash_InvalidVersion(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	err := enc.Matches("$argon2id$v=invalid$m=19456,t=2,p=1$AAAAAAAAAAAAAAAAAAAAAA==$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", "secret")
	if err == nil {
		t.Fatal("expected error for invalid version")
	}
	if err != ErrorInvalidHash {
		t.Fatalf("expected ErrorInvalidHash, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_InvalidHash_InvalidParams(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	err := enc.Matches("$argon2id$v=19$m=invalid,t=2,p=1$AAAAAAAAAAAAAAAAAAAAAA==$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", "secret")
	if err == nil {
		t.Fatal("expected error for invalid params")
	}
	if err != ErrorInvalidHash {
		t.Fatalf("expected ErrorInvalidHash, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_InvalidHash_InvalidSaltBase64(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	err := enc.Matches("$argon2id$v=19$m=19456,t=2,p=1$!!!invalid-base64!!!$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", "secret")
	if err == nil {
		t.Fatal("expected error for invalid salt base64")
	}
	if err != ErrorInvalidHash {
		t.Fatalf("expected ErrorInvalidHash, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_InvalidHash_InvalidHashBase64(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	// parts[4]=salt (valid), parts[5]=hash (invalid base64) - triggers DecodeString(parts[5]) error
	err := enc.Matches("$argon2id$v=19$m=19456,t=2,p=1$AAAAAAAAAAAAAAAAAAAAAA==$!!!invalid-base64!!!", "secret")
	if err == nil {
		t.Fatal("expected error for invalid hash base64")
	}
	if err != ErrorInvalidHash {
		t.Fatalf("expected ErrorInvalidHash, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_IncompatibleVersion(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	// Create valid hash first, then modify version in the string
	hash, err := enc.Encode("secret")
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}
	// Replace v=19 with v=18 to simulate incompatible version
	modifiedHash := strings.Replace(hash, "v=19", "v=18", 1)

	err = enc.Matches(modifiedHash, "secret")
	if err == nil {
		t.Fatal("expected error for incompatible version")
	}
	if err != ErrorIncompatibleVersion {
		t.Fatalf("expected ErrorIncompatibleVersion, got %v", err)
	}
}

func TestArgon2idPasswordEncoder_Mismatch_ReturnsCorrectError(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()

	hash, err := enc.Encode("super-secret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = enc.Matches(hash, "wrong-password")
	if err == nil {
		t.Fatal("expected mismatch error, got nil")
	}
	if err != ErrorPasswordMismatch {
		t.Fatalf("expected ErrorPasswordMismatch, got %v", err)
	}
}

func TestNewArgon2IDPasswordEncoder(t *testing.T) {
	enc := NewArgon2IDPasswordEncoder()
	if enc == nil {
		t.Fatal("expected non-nil encoder")
	}
}

func TestNewArgon2IDPasswordEncoder_NilReader(t *testing.T) {
	enc := newArgon2IDPasswordEncoder(nil)
	if enc == nil {
		t.Fatal("expected non-nil encoder when reader is nil")
	}
	// Should work - nil reader falls back to rand.Reader
	hash, err := enc.Encode("password")
	if err != nil {
		t.Fatalf("expected no error with nil reader fallback, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestArgon2idPasswordEncoder_Encode_RandReadError(t *testing.T) {
	errReader := &errReader{err: errors.New("rand read failed")}
	enc := newArgon2IDPasswordEncoder(errReader)

	_, err := enc.Encode("password")
	if err == nil {
		t.Fatal("expected error when rand reader fails")
	}
	if err.Error() != "rand read failed" {
		t.Fatalf("expected rand read error, got %v", err)
	}
}

type errReader struct{ err error }

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
