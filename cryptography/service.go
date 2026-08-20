package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/mlkem"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/sha3"
)

// JubJubPrimeGroup is the BabyJubJub base field modulus. Shared secrets are
// reduced mod this prime by the relayer so they fit in a Poseidon/BN254 field
// element ("salt"); the listener must reduce identically or HKDF inputs differ.
var JubJubPrimeGroup, _ = new(big.Int).SetString(
	"21888242871839275222246405745257275088548364400416034343698204186575808495617", 10,
)

// ReduceToSalt mirrors the relayer's GenerateSalt/RecoverSalt reduction
// (rayls-relayer, package cryptography): it interprets the raw ML-KEM shared
// secret as a big.Int, reduces it mod JubJubPrimeGroup, and returns the
// minimal big-endian byte representation. Output length is variable (1–32
// bytes) because big.Int.Bytes strips leading zeros — this matches the
// encryptor exactly. Any divergence in the reduction or the byte-encoding
// here silently desyncs HKDF and every DVP swap decrypts to garbage, so
// changes must be reviewed alongside the relayer counterpart.
func ReduceToSalt(sharedSecret []byte) []byte {
	salt := new(big.Int).SetBytes(sharedSecret)
	salt.Mod(salt, JubJubPrimeGroup)
	return salt.Bytes()
}

// ImportDecapsulationKey parses a hex-encoded ML-KEM 768 decapsulation key
func ImportDecapsulationKey(hexKey string) (*mlkem.DecapsulationKey768, error) {
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex key: %w", err)
	}

	dk, err := mlkem.NewDecapsulationKey768(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ML-KEM decapsulation key: %w", err)
	}

	return dk, nil
}

// Decapsulate recovers a shared secret from an ML-KEM ciphertext using a decapsulation key
func Decapsulate(dk *mlkem.DecapsulationKey768, ciphertext []byte) ([]byte, error) {
	secret, err := dk.Decapsulate(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("ML-KEM decapsulation failed: %w", err)
	}

	return secret, nil
}

// DeriveSymmetricKey derives a 256-bit symmetric key from a shared secret using HKDF-SHA3-256
func DeriveSymmetricKey(sharedSecret []byte) ([]byte, error) {
	context := []byte("Rayls")

	reader := hkdf.New(sha3.New256, sharedSecret, nil, context)

	key := make([]byte, 32)
	if _, err := io.ReadFull(reader, key); err != nil {
		return nil, fmt.Errorf("failed to derive symmetric key: %w", err)
	}

	return key, nil
}

// KMAC computes a keyed hash (message authentication code) using SHA3-256
func KMAC(key []byte, data []byte) []byte {
	h := sha3.New256()
	h.Write(key)
	h.Write(data)

	return h.Sum(nil)
}

// DecryptAuditPrivateKey decrypts a participant's encrypted private key using a shared secret.
// The shared secret is recovered via ML-KEM decapsulation of a key agreement ciphertext.
// The MAC is computed over the raw shared secret (matching the relayer's EncryptPrivateKey).
func DecryptAuditPrivateKey(encryptedPrivateKey []byte, mac []byte, sharedSecret []byte) ([]byte, error) {
	expectedMAC := KMAC(sharedSecret, encryptedPrivateKey)
	if !bytes.Equal(mac, expectedMAC) {
		return nil, fmt.Errorf("MAC verification failed")
	}

	symKey, err := DeriveSymmetricKey(sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to derive symmetric key: %w", err)
	}

	_, plaintext, err := DecryptGCM(encryptedPrivateKey, symKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// The relayer encrypts using gcmEncryptAndMarshal which json.Marshal's the []byte
	// before encryption. json.Marshal encodes []byte as a base64 JSON string,
	// so we must json.Unmarshal to recover the original raw bytes.
	var rawKey []byte
	if err := json.Unmarshal(plaintext, &rawKey); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted key: %w", err)
	}

	return rawKey, nil
}

// DecryptGCM receives a ciphertext and returns the plaintext (and potentially the associated data)
// this ciphertext has the following form: AD (16 bytes) || Nonce (12 bytes) || Encrypted Msg (remainder bytes)
func DecryptGCM(ciphertext []byte, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	associatedData := ciphertext[0:16]
	nonce := ciphertext[16 : 16+gcm.NonceSize()]
	ctxt := ciphertext[16+gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ctxt, associatedData)
	if err != nil {
		return nil, nil, err
	}
	return associatedData, plaintext, nil
}

// HashIt hashes data using SHA3-256
func HashIt(data []byte) []byte {
	h := sha3.New256()
	h.Write(data)

	return h.Sum(nil)
}
