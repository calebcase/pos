package aesprng

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/calebcase/pos"
)

type SeedSizeError int

func (e SeedSizeError) Error() string {
	return fmt.Sprintf("Invalid seed size %d", int(e))
}

type State struct {
	key []byte
	iv  []byte

	mode cipher.BlockMode
	zero []byte
}

var _ pos.PRNG = (*State)(nil)

func New(key, iv []byte) (prng *State, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &State{
		key: append([]byte(nil), key...),
		iv:  append([]byte(nil), iv...),

		mode: cipher.NewCBCEncrypter(block, iv),
		zero: make([]byte, 1024, 1024),
	}, nil
}

func (prng *State) Read(b []byte) (n int, err error) {
	if len(prng.zero) < len(b) {
		prng.zero = make([]byte, len(b), len(b))
	}

	prng.mode.CryptBlocks(b, prng.zero[:len(b)])

	return len(b), nil
}

func (prng *State) New(seed []byte) (nprng pos.PRNG, err error) {
	var offset int
	switch len(seed) {
	// AES-256
	case 32 + 16:
		offset = 32
	// AES-192
	case 24 + 16:
		offset = 24
	// AES-128
	case 16 + 16:
		offset = 16
	default:
		return nil, SeedSizeError(len(seed))
	}

	key := seed[:offset]
	iv := seed[offset:]

	return New(key, iv)
}

func (prng *State) Clone() (nprng pos.PRNG, err error) {
	return New(prng.key, prng.iv)
}
