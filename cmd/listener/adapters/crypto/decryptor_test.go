package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/mlkem"
	"crypto/rand"
	"encoding/json"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/testutil"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/config"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cryptography"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/types"
)

// setupMLKEMSharedSecret generates an ML-KEM key pair, encapsulates to produce a shared secret
func setupMLKEMSharedSecret(t *testing.T) []byte {
	t.Helper()
	dk, err := mlkem.GenerateKey768()
	require.NoError(t, err)
	ek := dk.EncapsulationKey()
	ss, _ := ek.Encapsulate()
	return ss
}

// encryptGCMWithDeriveKey encrypts plaintext using DeriveSymmetricKey + AES-GCM (AD || Nonce || Ciphertext format)
func encryptGCMWithDeriveKey(t *testing.T, plaintext []byte, sharedSecret []byte) []byte {
	t.Helper()
	symKey, err := cryptography.DeriveSymmetricKey(sharedSecret)
	require.NoError(t, err)
	ct, err := encryptGCM(plaintext, symKey)
	require.NoError(t, err)
	return ct
}

// encryptGCM is a helper that produces ciphertext in AD || Nonce || EncryptedMsg format
func encryptGCM(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	associatedData := make([]byte, 16)
	if _, err = io.ReadFull(rand.Reader, associatedData); err != nil {
		return nil, err
	}

	encryptedMsg := gcm.Seal(nil, nonce, plaintext, associatedData)

	// Format: AD || Nonce || Encrypted Msg
	ciphertext := make([]byte, 0, len(associatedData)+len(nonce)+len(encryptedMsg))
	ciphertext = append(ciphertext, associatedData...)
	ciphertext = append(ciphertext, nonce...)
	ciphertext = append(ciphertext, encryptedMsg...)

	return ciphertext, nil
}

// TestDecryptor covers core behaviors of the decryptor adapter
func TestDecryptor(t *testing.T) {
	// Generate two ML-KEM shared secrets (simulating two participants)
	ssAB := setupMLKEMSharedSecret(t)  // shared secret between participants A and B
	ssVen := setupMLKEMSharedSecret(t) // Hub shared secret

	type testMessage struct {
		Message string `json:"message"`
		Value   int    `json:"value"`
	}
	originalMessage := testMessage{Message: "hello world", Value: 42}
	plaintext, err := json.Marshal(originalMessage)
	require.NoError(t, err)

	stubLogger := &testutil.StubLogger{}
	decryptor := NewDecryptor(nil, &config.Config{}, stubLogger)

	// Encrypt with A-B shared secret for participant tests
	encryptedForAB := encryptGCMWithDeriveKey(t, plaintext, ssAB)
	// Encrypt with Hub shared secret for atomic tests
	encryptedForVen := encryptGCMWithDeriveKey(t, plaintext, ssVen)

	buildParticipantPL := func() core.PNodeDataAndSecrets {
		return core.PNodeDataAndSecrets{
			// Pairwise shared secret entry
			"1-2": {
				"100": {
					SharedSecret: ssAB,
					ChainId:      big.NewInt(2),
					BlockNumber:  big.NewInt(100),
				},
			},
			// Hub entries (also included in participant search)
			"1": {
				"100": {
					HubSharedSecret: ssVen,
					ChainId:         big.NewInt(1),
					BlockNumber:     big.NewInt(100),
				},
			},
		}
	}

	buildAtomicPL := func(ss []byte) core.PNodeDataAndSecrets {
		return core.PNodeDataAndSecrets{
			"2": {"200": {HubSharedSecret: ss, ChainId: big.NewInt(2), BlockNumber: big.NewInt(200)}},
		}
	}

	type testCase struct {
		name        string
		secretType  types.SecretType
		buildPL     func(ss []byte) core.PNodeDataAndSecrets
		payload     []byte
		blockNumber uint64
		expectErr   string
		assertFn    func(t *testing.T, got []byte)
	}

	testCases := []testCase{
		{
			name:        "participant secret success",
			secretType:  types.ParticipantSecret,
			buildPL:     func(_ []byte) core.PNodeDataAndSecrets { return buildParticipantPL() },
			payload:     encryptedForAB,
			blockNumber: 100,
			assertFn: func(t *testing.T, got []byte) {
				var msg testMessage
				require.NoError(t, json.Unmarshal(got, &msg))
				assert.Equal(t, originalMessage, msg)
			},
		},
		{
			name:        "atomic secret success",
			secretType:  types.AtomicSecret,
			buildPL:     func(ss []byte) core.PNodeDataAndSecrets { return buildAtomicPL(ss) },
			payload:     encryptedForVen,
			blockNumber: 200,
			assertFn: func(t *testing.T, got []byte) {
				var msg testMessage
				require.NoError(t, json.Unmarshal(got, &msg))
				assert.Equal(t, originalMessage, msg)
			},
		},
		{
			name:       "error - no key",
			secretType: types.ParticipantSecret,
			buildPL: func(_ []byte) core.PNodeDataAndSecrets {
				wrongSS := make([]byte, 32)
				return core.PNodeDataAndSecrets{
					"1": {
						"100": {
							HubSharedSecret: wrongSS,
							ChainId:         big.NewInt(1),
							BlockNumber:     big.NewInt(100),
						},
					},
				}
			},
			payload:     encryptedForAB,
			blockNumber: 100,
			expectErr:   "no key could decrypt the payload",
		},
		{
			name:        "error - unknown secret type",
			secretType:  types.SecretType(77),
			buildPL:     func(_ []byte) core.PNodeDataAndSecrets { return nil },
			payload:     encryptedForAB,
			blockNumber: 10,
			expectErr:   "unrecognized secretType",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pnData := tc.buildPL(ssVen)
			got, err := decryptor.DecryptPayloadBytes(tc.payload, tc.blockNumber, pnData, tc.secretType)
			if tc.expectErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErr)
				return
			}
			require.NoError(t, err)
			if tc.assertFn != nil {
				tc.assertFn(t, got)
			}
		})
	}

	t.Run("filter valid selects latest <= block", func(t *testing.T) {
		ss1 := []byte{1, 2, 3}
		pnData := core.PNodeDataAndSecrets{"1": {
			"90":  {SharedSecret: ss1, ChainId: big.NewInt(1), BlockNumber: big.NewInt(90)},
			"110": {SharedSecret: ss1, ChainId: big.NewInt(1), BlockNumber: big.NewInt(110)},
		}}
		filtered := filterValidDataAndSecrets(pnData, 105)
		_, exists90 := filtered["1"]["90"]
		_, exists110 := filtered["1"]["110"]
		assert.True(t, exists90)
		assert.False(t, exists110)
	})
}

// buildParticipantDK generates a participant ML-KEM keypair, encapsulates
// against it to produce a (ctxt, sharedSecret) pair, and returns everything
// needed to simulate a SwapInitiated ciphertext arriving at the listener.
func buildParticipantDK(t *testing.T) (dkBytes []byte, ctxt []byte, salt []byte) {
	t.Helper()
	dk, err := mlkem.GenerateKey768()
	require.NoError(t, err)
	ek := dk.EncapsulationKey()
	ss, ct := ek.Encapsulate()
	return dk.Bytes(), ct, cryptography.ReduceToSalt(ss)
}

// TestDecryptor_DecryptWithSalt_RoundTrip exercises the SwapCompleted path:
// encrypt with a caller-supplied salt, decrypt with the same salt.
func TestDecryptor_DecryptWithSalt_RoundTrip(t *testing.T) {
	// DecryptWithSalt returns the plaintext unchanged when the salt matches.
	decryptor := NewDecryptor(nil, &config.Config{}, &testutil.StubLogger{})
	salt := []byte("a-reduced-salt-of-some-length")
	plaintext := []byte(`{"hello":"world"}`)
	ciphertext := encryptGCMWithDeriveKey(t, plaintext, salt)

	got, err := decryptor.DecryptWithSalt(ciphertext, salt)

	require.NoError(t, err)
	assert.Equal(t, plaintext, got)
}

// TestDecryptor_DecryptWithSalt_WrongSalt ensures a mismatched salt fails cleanly.
func TestDecryptor_DecryptWithSalt_WrongSalt(t *testing.T) {
	// A salt that doesn't match the one used at encryption time must fail, not panic.
	decryptor := NewDecryptor(nil, &config.Config{}, &testutil.StubLogger{})
	ciphertext := encryptGCMWithDeriveKey(t, []byte("payload"), []byte("right-salt"))

	_, err := decryptor.DecryptWithSalt(ciphertext, []byte("wrong-salt"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt payload")
}

// TestDecryptor_DecryptSwapPayload_FindsMatchingParticipant proves the
// candidate-DK loop selects the right participant among several.
func TestDecryptor_DecryptSwapPayload_FindsMatchingParticipant(t *testing.T) {
	// With two participant DKs recorded, the decryptor must try each until
	// decapsulation yields a salt that decrypts the payload.
	targetDK, targetCtxt, targetSalt := buildParticipantDK(t)
	decoyDK, _, _ := buildParticipantDK(t)

	plaintext := []byte(`{"sharedId":"abc"}`)
	encrypted := encryptGCMWithDeriveKey(t, plaintext, targetSalt)

	pnData := core.PNodeDataAndSecrets{
		"1": {"100": {ChainId: big.NewInt(1), BlockNumber: big.NewInt(100), ParticipantDK: decoyDK}},
		"2": {"100": {ChainId: big.NewInt(2), BlockNumber: big.NewInt(100), ParticipantDK: targetDK}},
	}

	decryptor := NewDecryptor(nil, &config.Config{}, &testutil.StubLogger{})
	got, gotSalt, err := decryptor.DecryptSwapPayload(targetCtxt, encrypted, 100, pnData)

	require.NoError(t, err)
	assert.Equal(t, plaintext, got)
	assert.Equal(t, targetSalt, gotSalt)
}

// TestDecryptor_DecryptSwapPayload_NoCandidate fails cleanly when no DK can
// decapsulate into a salt that decrypts the payload.
func TestDecryptor_DecryptSwapPayload_NoCandidate(t *testing.T) {
	// With only a decoy DK present, no candidate can derive a salt matching
	// the target ciphertext — the decryptor must return an error, not panic.
	_, targetCtxt, targetSalt := buildParticipantDK(t)
	decoyDK, _, _ := buildParticipantDK(t)

	encrypted := encryptGCMWithDeriveKey(t, []byte("payload"), targetSalt)

	pnData := core.PNodeDataAndSecrets{
		"1": {"100": {ChainId: big.NewInt(1), BlockNumber: big.NewInt(100), ParticipantDK: decoyDK}},
	}

	decryptor := NewDecryptor(nil, &config.Config{}, &testutil.StubLogger{})
	_, _, err := decryptor.DecryptSwapPayload(targetCtxt, encrypted, 100, pnData)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no participant DK could decrypt swap payload")
}

// TestDecryptor_DecryptSwapPayload_SkipsEntriesWithoutDK verifies entries
// missing ParticipantDK (e.g. pairwise-secret-only rows) are silently skipped.
func TestDecryptor_DecryptSwapPayload_SkipsEntriesWithoutDK(t *testing.T) {
	// A pairwise-secret-only entry (no ParticipantDK) must not short-circuit
	// the loop — the decryptor should move on to the next candidate.
	targetDK, targetCtxt, targetSalt := buildParticipantDK(t)
	encrypted := encryptGCMWithDeriveKey(t, []byte("payload"), targetSalt)

	pnData := core.PNodeDataAndSecrets{
		"1-2": {"100": {ChainId: big.NewInt(2), BlockNumber: big.NewInt(100), SharedSecret: []byte{1, 2, 3}}},
		"2":   {"100": {ChainId: big.NewInt(2), BlockNumber: big.NewInt(100), ParticipantDK: targetDK}},
	}

	decryptor := NewDecryptor(nil, &config.Config{}, &testutil.StubLogger{})
	got, _, err := decryptor.DecryptSwapPayload(targetCtxt, encrypted, 100, pnData)

	require.NoError(t, err)
	assert.Equal(t, []byte("payload"), got)
}

// TestReduceToSalt_IsDeterministic proves identical input yields identical
// output — the listener and relayer must agree on this reduction.
func TestReduceToSalt_IsDeterministic(t *testing.T) {
	// Two calls with the same raw secret must produce byte-identical salts,
	// otherwise HKDF diverges between encryptor and decryptor.
	secret := make([]byte, 32)
	for i := range secret {
		secret[i] = byte(i)
	}

	got1 := cryptography.ReduceToSalt(secret)
	got2 := cryptography.ReduceToSalt(secret)

	assert.Equal(t, got1, got2)
}

// TestReduceToSalt_ReducesModPrime confirms inputs >= JubJubPrimeGroup are
// reduced (not passed through) and the returned bytes reflect the modulus.
func TestReduceToSalt_ReducesModPrime(t *testing.T) {
	// Feed in exactly JubJubPrimeGroup: the result must be 0 (empty bytes
	// under big.Int.Bytes semantics), proving the Mod is applied.
	primeBytes := cryptography.JubJubPrimeGroup.Bytes()

	got := cryptography.ReduceToSalt(primeBytes)

	assert.Empty(t, got, "prime mod prime should be zero, encoded as empty bytes")
}

// TestReduceToSalt_StripsLeadingZeros documents the variable-length output
// contract — callers that need fixed-width salts must left-pad themselves.
func TestReduceToSalt_StripsLeadingZeros(t *testing.T) {
	// big.Int.Bytes strips leading zeros, so a small input reduces to a
	// short output. The listener and relayer both rely on this semantic.
	small := []byte{0, 0, 0, 0, 0, 0, 0, 42}

	got := cryptography.ReduceToSalt(small)

	assert.Equal(t, []byte{42}, got)
}
