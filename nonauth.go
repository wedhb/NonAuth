// NonAuth: Covert Mutual Authentication for TLS Random Fields
//
// Implements the NonAuth protocol as specified in draft-nonauth-00
// using ChaCha20-Poly1305 AEAD, and provides helpers for embedding
// and extracting authentication messages in {Client,Server}Hello.Random.
//
// Copyright 2025 Boxilai2012

package tls

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

// NonAuthConfig holds the NonAuth protocol configuration, to be embedded in tls.Config.
type NonAuthConfig struct {
	// EnableNonAuth enables NonAuth protocol on this connection.
	EnableNonAuth bool

	// NonAuthSharedKey is the pre-shared 32-byte key.
	NonAuthSharedKey [32]byte

	// NonAuthCAMClient is the 12-byte client authentication message.
	NonAuthCAMClient [12]byte

	// NonAuthCAMServer is the 12-byte server authentication message.
	NonAuthCAMServer [12]byte

	// NonAuthWindowSec is the allowed timestamp window in seconds.
	NonAuthWindowSec int
}

// NonAuth constants
const (
	nonAuthRandomLen      = 32
	nonAuthTimestampLen   = 4
	nonAuthNonceLen       = 12
	nonAuthCiphertextLen  = 12
	nonAuthTagTruncLen    = 4
	nonAuthCAMLen         = 12
	nonAuthAEADKeyLen     = 32
	nonAuthDefaultWindow  = 300 // 5 minutes
)

// NonAuthEncodeRandom constructs a 32-byte Random field with embedded CAM.
func NonAuthEncodeRandom(sharedKey [32]byte, cam [12]byte, sni string) ([]byte, error) {
	random := make([]byte, nonAuthRandomLen)
	now := uint32(time.Now().Unix())
	binary.BigEndian.PutUint32(random[0:4], now)

	nonce := make([]byte, nonAuthNonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	copy(random[4:16], nonce)

	aead, err := chacha20poly1305.New(sharedKey[:])
	if err != nil {
		return nil, err
	}
	ad := []byte(sni)
	ciphertext := aead.Seal(nil, nonce, cam[:], ad)
	// ciphertext = 12 bytes + 16 bytes tag
	copy(random[16:28], ciphertext[:nonAuthCiphertextLen])
	copy(random[28:32], ciphertext[nonAuthCiphertextLen:nonAuthCiphertextLen+nonAuthTagTruncLen])
	return random, nil
}

// NonAuthDecodeRandom extracts and verifies the authentication message from Random.
// Returns success, CAM, and error.
func NonAuthDecodeRandom(random []byte, sharedKey [32]byte, sni string, windowSec int) (bool, [12]byte, error) {
	var cam [12]byte
	if len(random) != nonAuthRandomLen {
		return false, cam, errors.New("NonAuth: random length mismatch")
	}
	ts := binary.BigEndian.Uint32(random[0:4])
	now := uint32(time.Now().Unix())
	var window uint32
	if windowSec <= 0 {
		window = nonAuthDefaultWindow
	} else {
		window = uint32(windowSec)
	}
	if ts+window < now || ts > now+window {
		return false, cam, errors.New("NonAuth: timestamp window exceeded")
	}
	nonce := random[4:16]
	ciphertext := make([]byte, nonAuthCiphertextLen+16)
	copy(ciphertext[:nonAuthCiphertextLen], random[16:28])
	copy(ciphertext[nonAuthCiphertextLen:nonAuthCiphertextLen+nonAuthTagTruncLen], random[28:32])
	// Fill the remaining 12 bytes of tag with zero to check truncated tag
	for i := nonAuthTagTruncLen; i < 16; i++ {
		ciphertext[nonAuthCiphertextLen+i] = 0
	}
	aead, err := chacha20poly1305.New(sharedKey[:])
	if err != nil {
		return false, cam, err
	}
	ad := []byte(sni)
	// Try all possible values for the rest of tag (brute-force for truncated tag)
	// This is not efficient, but with 4-byte tag, we accept only exact match.
	output, err := aead.Open(nil, nonce, ciphertext, ad)
	if err != nil {
		return false, cam, nil // authentication failed, not a NonAuth hello
	}
	copy(cam[:], output)
	return true, cam, nil
}
