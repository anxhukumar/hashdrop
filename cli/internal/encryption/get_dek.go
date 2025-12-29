package encryption

import (
	"github.com/anxhukumar/hashdrop/cli/internal/config"
	"golang.org/x/crypto/argon2"
)

func DeriveDEK(passphrase string, salt []byte) []byte {
	passphraseBytes := []byte(passphrase)

	key := argon2.IDKey(
		passphraseBytes,
		salt,
		config.ArgonTime,
		config.ArgonMemory,
		config.ArgonThreads,
		config.ArgonKeyLen,
	)

	return key
}
