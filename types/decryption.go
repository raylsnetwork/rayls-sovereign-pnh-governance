package types

import "math/big"

// IPNodeDataAndSecrets contains pre-computed shared secrets for Privacy Node data decryption.
// SharedSecret is used for participant-to-participant decryption (ParticipantSecret),
// while HubSharedSecret is used for operator-participant decryption (AtomicSecret).
// ParticipantDK holds the participant's raw ML-KEM decapsulation key bytes so events
// whose ciphertext is encapsulated against a participant's view key (e.g. SwapInitiated)
// can be decrypted by trying each candidate.
type IPNodeDataAndSecrets struct {
	ChainId         *big.Int
	BlockNumber     *big.Int
	SharedSecret    []byte // pairwise shared secret from ML-KEM key agreement (32 bytes)
	HubSharedSecret []byte // shared secret with Private Network Hub operator from ML-KEM key agreement (32 bytes)
	ParticipantDK   []byte // participant's ML-KEM decapsulation key (raw bytes, reconstructable via mlkem.NewDecapsulationKey768)
}

// SecretType represents the type of secret used for decryption
type SecretType int

const (
	ParticipantSecret SecretType = iota
	AtomicSecret
)
