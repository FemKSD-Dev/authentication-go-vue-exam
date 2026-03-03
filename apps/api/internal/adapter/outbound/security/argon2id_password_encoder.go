package security

import (
	"authentication-project-exam/internal/core/port/outbound"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrorInvalidHash         = errors.New("invalid hash")
	ErrorIncompatibleVersion = errors.New("incompatible argon2 version")
	ErrorPasswordMismatch    = errors.New("password mismatch")
)

type Argon2IDPasswordEncoder struct {
	memory     uint32
	time       uint32
	threads    uint8
	keyLength  uint32
	saltLength uint32
	randReader io.Reader
}

type argon2Params struct {
	memory  uint32
	time    uint32
	threads uint8
	version int
}

func NewArgon2IDPasswordEncoder() outbound.PasswordManager {
	return newArgon2IDPasswordEncoder(rand.Reader)
}

func newArgon2IDPasswordEncoder(randReader io.Reader) *Argon2IDPasswordEncoder {
	if randReader == nil {
		randReader = rand.Reader
	}
	return &Argon2IDPasswordEncoder{
		memory:     19 * 1024,
		time:       2,
		threads:    1,
		keyLength:  32,
		saltLength: 16,
		randReader: randReader,
	}
}

func (e *Argon2IDPasswordEncoder) Encode(raw string) (string, error) {
	salt := make([]byte, e.saltLength)
	if _, err := io.ReadFull(e.randReader, salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(raw), salt, e.time, e.memory, e.threads, e.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, e.memory, e.time, e.threads, b64Salt, b64Hash)
	return encoded, nil
}

func (e *Argon2IDPasswordEncoder) Matches(encoded, raw string) error {
	params, salt, expectedHash, err := decodeArgon2IDHash(encoded)
	if err != nil {
		return err
	}
	if params.version != argon2.Version {
		return ErrorIncompatibleVersion
	}

	actualHash := argon2.IDKey([]byte(raw), salt, params.time, params.memory, params.threads, uint32(len(expectedHash)))

	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return ErrorPasswordMismatch
	}
	return nil
}

func decodeArgon2IDHash(encoded string) (*argon2Params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrorInvalidHash
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, ErrorInvalidHash
	}

	params := &argon2Params{}

	if _, err := fmt.Sscanf(parts[2], "v=%d", &params.version); err != nil {
		return nil, nil, nil, ErrorInvalidHash
	}

	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.memory, &params.time, &params.threads); err != nil {
		return nil, nil, nil, ErrorInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, ErrorInvalidHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, ErrorInvalidHash
	}

	return params, salt, hash, nil
}
