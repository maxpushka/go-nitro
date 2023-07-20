package crypto

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/statechannels/go-nitro/types"
)

// Signature is an ECDSA signature
type Signature struct {
	R []byte
	S []byte
	V byte
}

// SignEthereumMessage accepts an arbitrary message, prepends a known message,
// hashes the result using keccak256 and calculates the secp256k1 signature
// of the hash using the provided secret key. The known message added to the input before hashing is
// "\x19Ethereum Signed Message:\n" + len(message).
// See https://github.com/ethereum/go-ethereum/pull/2940 and EIPs 191, 721.
func SignEthereumMessage(message []byte, secretKey []byte) (Signature, error) {
	digest := computeEthereumSignedMessageDigest(message)
	concatenatedSignature, error := secp256k1.Sign(digest, secretKey)
	if error != nil {
		return Signature{}, error
	}
	sig := SplitSignature(concatenatedSignature)

	// This step is necessary to remain compatible with the ecrecover precompile
	if int(sig.V) < 27 {
		sig.V = byte(int(sig.V + 27))
	}

	return sig, nil
}

// RecoverEthereumMessageSigner accepts a message (bytestring) and signature generated by SignEthereumMessage.
// It reconstructs the appropriate digest and recovers an address via secp256k1 public key recovery
func RecoverEthereumMessageSigner(message []byte, signature Signature) (common.Address, error) {
	// This step is necessary to remain compatible with the ecrecover precompile
	sig := signature
	if int(sig.V) >= 27 {
		sig.V = byte(int(sig.V - 27))
	}

	digest := computeEthereumSignedMessageDigest(message)
	pubKey, error := secp256k1.RecoverPubkey(digest, joinSignature(sig))
	if error != nil {
		return types.Address{}, error
	}
	ecdsaPubKey, error := crypto.UnmarshalPubkey(pubKey)
	if error != nil {
		return types.Address{}, error
	}
	crypto.PubkeyToAddress(*ecdsaPubKey)
	return crypto.PubkeyToAddress(*ecdsaPubKey), error
}

// computeEthereumSignedMessageDigest accepts an arbitrary message, prepends a known message,
// and hashes the result using keccak256. The known message added to the input before hashing is
// "\x19Ethereum Signed Message:\n" + len(message).
func computeEthereumSignedMessageDigest(message []byte) []byte {
	return crypto.Keccak256(
		[]byte(
			fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), string(message)),
		),
	)
}

// splitSignature takes a 65 bytes signature in the [R||S||V] format and returns the individual components
func SplitSignature(concatenatedSignature []byte) (signature Signature) {
	signature.R = concatenatedSignature[:32]
	signature.S = concatenatedSignature[32:64]
	signature.V = concatenatedSignature[64]
	return
}

// joinSignature takes a Signature and returns the concatenatedSignature in the [R||S||V] format
func joinSignature(signature Signature) (concatenatedSignature []byte) {
	concatenatedSignature = append(concatenatedSignature, signature.R...)
	concatenatedSignature = append(concatenatedSignature, signature.S...)
	concatenatedSignature = append(concatenatedSignature, signature.V)
	return
}

// ToHexString returns the signature as a hex string
func (s Signature) ToHexString() string {
	return hexutil.Encode(joinSignature(s))
}

func (s1 Signature) Equal(s2 Signature) bool {
	return bytes.Equal(s1.S, s2.S) && bytes.Equal(s1.R, s2.R) && s1.V == s2.V
}

func (s Signature) MarshalJSON() ([]byte, error) {
	joined := joinSignature(s)
	hex := hexutil.Encode(joined)
	return json.Marshal(hex)
}

func (s *Signature) UnmarshalJSON(b []byte) error {
	var hex string
	err := json.Unmarshal(b, &hex)
	if err != nil {
		return err
	}
	joined, err := hexutil.Decode(hex)
	if err != nil {
		return err
	}

	// If the signature is all zeros, we consider it to be the empty signature
	if allZero(joined) {
		return nil
	}

	if len(joined) != 65 {
		return fmt.Errorf("signature must be 65 bytes long or a zero string, received %d bytes", len(joined))
	}

	s.R = joined[:32]
	s.S = joined[32:64]
	s.V = joined[64]
	return nil
}

// allZero returns true if all bytes in the slice are zero false otherwise
func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}
