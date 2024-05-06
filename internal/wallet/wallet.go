package wallet

import (
	"GoBlockchain/internal/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	// ECDSA - Elliptic Curve Digital Signature Algorithm
	PrivateKey SerializablePrivateKey
	PublicKey  []byte
}

type SerializablePrivateKey struct {
	D         *big.Int
	PublicKey ecdsa.PublicKey
}

func NewKeyPair() (SerializablePrivateKey, []byte) {
	curve := elliptic.P256()

	// Use the curve with a random number generator to generate a unique private key
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	// To make the public key, we append the X and Y components of the private key
	// These points are coordinates on the elliptic curve which we used to generate the private key
	// This is secure because it is computationally infeasible to reverse-engineer the private key from the coordinates
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	// Create a serializable private key
	serializablePrivateKey := SerializablePrivateKey{
		D: private.D,
		PublicKey: ecdsa.PublicKey{
			X: private.PublicKey.X,
			Y: private.PublicKey.Y,
		},
	}

	return serializablePrivateKey, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	Wallet := Wallet{private, public}

	return &Wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)
	return pubHash[:]
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func (w Wallet) Address() []byte {
	/** In order to create an address we do the following:
	1. Get the public key hash
	2. Add the version byte in front of the public key hash
	3. Get the checksum of the versioned public key hash
	4. Append the checksum to the versioned public key hash
	5. Encode the result into base58
	*/
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := utils.Base58Encode(fullHash)

	return address
}

func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}
