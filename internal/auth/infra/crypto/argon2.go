package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2Hasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewArgon2Hasher() *argon2Hasher {
	return &argon2Hasher{
		memory:      64 * 1024, // 64 MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

func (h *argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generando salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, h.iterations, h.memory, h.parallelism, h.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", h.memory, h.iterations, h.parallelism, b64Salt, b64Hash)
	return encoded, nil
}

func (h *argon2Hasher) Compare(encodedHash, password string) error {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return errors.New("Formato de hash invalido")
	}

	var memory, iterations uint32
	var parallelism uint8

	params := strings.Split(parts[3], ",")
	for _, param := range params {
		if strings.HasPrefix(param, "m=") {
			fmt.Sscanf(param, "m=%d", &memory)
		} else if strings.HasPrefix(param, "t=") {
			fmt.Sscanf(param, "t=%d", &iterations)
		} else if strings.HasPrefix(param, "p=") {
			fmt.Sscanf(param, "p=%d", &parallelism)
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("decodificando salt: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("decodificando hash: %w", err)
	}

	gotHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(hash)))

	if subtle.ConstantTimeCompare(gotHash, hash) != 1 {
		return errors.New("contraseña incorrecta")
	}
	return nil
}
