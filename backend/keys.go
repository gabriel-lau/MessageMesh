package backend

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"sort"

	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	bolt "go.etcd.io/bbolt"
)

type KeyPair struct {
	PrivKey    libp2pcrypto.PrivKey // Private Key with libp2p crypto
	PubKey     libp2pcrypto.PubKey  // Public Key with libp2p crypto
	PrivateKey crypto.PrivateKey    // Private Key with standard crypto
	PublicKey  crypto.PublicKey     // Public Key with standard crypto
}

const (
	directory = "db"
	file      = "keys.db"
	dbpath    = directory + "/" + file
)

// Generate a new key pair
func NewKeyPair() (KeyPair, error) {
	keypair := KeyPair{}

	// Create new directory (db)
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return keypair, err
	}
	prvkey, pubkey, err := libp2pcrypto.GenerateKeyPairWithReader(libp2pcrypto.RSA, 2048, rand.Reader)
	if err != nil {
		fmt.Println("Error generating key pair:", err)
		return keypair, err
	}

	// Store keys
	boltDB, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return keypair, err
	}
	defer boltDB.Close()

	// Store Private Key and Public Key
	marshalledPrivKey, err := libp2pcrypto.MarshalPrivateKey(prvkey)
	if err != nil {
		fmt.Println("Error getting raw private key:", err)
		return keypair, err
	}
	marshalledPubKey, err := libp2pcrypto.MarshalPublicKey(pubkey)
	if err != nil {
		fmt.Println("Error getting raw public key:", err)
		return keypair, err
	}

	err = boltDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("keys"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = bucket.Put([]byte("private"), marshalledPrivKey)
		if err != nil {
			return fmt.Errorf("put: %s", err)
		}
		err = bucket.Put([]byte("public"), marshalledPubKey)
		if err != nil {
			return fmt.Errorf("put: %s", err)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error writing private key to database:", err)
		return keypair, err
	}

	keypair.PrivKey = prvkey
	keypair.PubKey = pubkey

	keypair.PrivateKey, err = keypair.getStdPrivateKey(prvkey)
	if err != nil {
		fmt.Println("Error getting raw private key:", err)
		return keypair, err
	}

	keypair.PublicKey, err = keypair.getStdPublicKey(pubkey)
	if err != nil {
		fmt.Println("Error getting raw public key:", err)
		return keypair, err
	}
	return keypair, nil
}

// Get the KeyPair from the database
func ReadKeyPair() (KeyPair, error) {
	keyPair := KeyPair{}
	boltDB, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return keyPair, err
	}
	defer boltDB.Close()

	err = boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("keys"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		privKeyBytes := bucket.Get([]byte("private"))
		pubKeyBytes := bucket.Get([]byte("public"))

		privKey, err := libp2pcrypto.UnmarshalPrivateKey(privKeyBytes)
		if err != nil {
			fmt.Println("Error getting private key:", err)
			return err
		}
		pubKey, err := libp2pcrypto.UnmarshalPublicKey(pubKeyBytes)
		if err != nil {
			fmt.Println("Error getting public key:", err)
			return err
		}

		keyPair.PrivKey = privKey
		keyPair.PubKey = pubKey

		keyPair.PrivateKey, err = keyPair.getStdPrivateKey(privKey)
		if err != nil {
			fmt.Println("Error getting raw private key:", err)
			return err
		}

		keyPair.PublicKey, err = keyPair.getStdPublicKey(pubKey)
		if err != nil {
			fmt.Println("Error getting raw public key:", err)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error reading private key from database:", err)
		return keyPair, err
	}
	return keyPair, nil
}

func (keypair *KeyPair) getStdPrivateKey(privKey libp2pcrypto.PrivKey) (crypto.PrivateKey, error) {
	stdPrivKey, err := libp2pcrypto.PrivKeyToStdKey(privKey)
	if err != nil {
		fmt.Println("Error getting raw private key:", err)
		return "", err
	}
	return stdPrivKey, nil
}

func (keypair *KeyPair) getStdPublicKey(pubKey libp2pcrypto.PubKey) (crypto.PublicKey, error) {
	stdPubKey, err := libp2pcrypto.PubKeyToStdKey(pubKey)
	if err != nil {
		fmt.Println("Error getting raw public key:", err)
		return "", err
	}
	return stdPubKey, nil
}

func (keypair *KeyPair) SignWithPrivateKey(plaintext []byte) ([]byte, error) {
	ciphertext, err := keypair.PrivKey.Sign(plaintext)
	if err != nil {
		fmt.Println("Error encrypting message:", err)
		return nil, err
	}
	return ciphertext, nil
}

func (keypair *KeyPair) VerifyWithPublicKey(plaintext []byte, sig []byte) (bool, error) {
	verification, err := keypair.PubKey.Verify(plaintext, sig)
	if err != nil {
		fmt.Println("Error encrypting message:", err)
		return false, err
	}
	return verification, nil
}

func (keypair *KeyPair) EncryptWithPublicKey(plaintext []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, keypair.PublicKey.(*rsa.PublicKey), plaintext)
	if err != nil {
		fmt.Println("Error encrypting message:", err)
		return nil, err
	}
	return ciphertext, nil
}

func (keypair *KeyPair) DecryptWithPrivateKey(ciphertext []byte) ([]byte, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, keypair.PrivateKey.(*rsa.PrivateKey), ciphertext)
	if err != nil {
		fmt.Println("Error decrypting message:", err)
		return nil, err
	}
	return plaintext, nil
}

// GenerateSymmetricKey creates a random symmetric key of specified length
func GenerateSymmetricKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Error generating symmetric key:", err)
		return nil, err
	}
	return key, nil
}

// EncryptWithSymmetricKey encrypts a message using AES-GCM with the symmetric key
func EncryptWithSymmetricKey(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("Error generating nonce:", err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating GCM:", err)
		return nil, err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptWithSymmetricKey decrypts a message using AES-GCM with the symmetric key
func DecryptWithSymmetricKey(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Error creating GCM:", err)
		return nil, err
	}

	if len(ciphertext) < 12 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:12], ciphertext[12:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("Error decrypting message:", err)
		return nil, err
	}

	return plaintext, nil
}

func SaveSymmetricKey(key []byte, peerIDs []string) error {
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("symmetric"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		sort.Strings(peerIDs)
		return bucket.Put([]byte(peerIDs[0]+peerIDs[1]), key)
	})
}

func GetSymmetricKey(peerIDs []string) ([]byte, error) {
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return nil, err
	}
	var key []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("symmetric"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		sort.Strings(peerIDs)
		key = bucket.Get([]byte(peerIDs[0] + peerIDs[1]))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, fmt.Errorf("symmetric key not found")
	}
	return key, nil
}

// GetPeerPublicKey retrieves the public key of a peer from their peer ID
func GetPeerPublicKey(p2p *P2PService, peerIDStr string) (libp2pcrypto.PubKey, error) {
	// Parse the peer ID string
	peerID, err := peer.Decode(peerIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode peer ID: %s", err.Error())
	}

	// Extract public key from peer ID
	pubKey, err := peerID.ExtractPublicKey()
	if err != nil {
		// If the public key isn't embedded in the peer ID, try to get it from the peer store
		pubKey = p2p.Host.Peerstore().PubKey(peerID)
		if pubKey == nil {
			return nil, fmt.Errorf("couldn't find public key for peer %s", peerIDStr)
		}
	}

	return pubKey, nil
}

// GetPeerStandardPublicKey retrieves the public key of a peer in standard crypto.PublicKey format
func GetPeerStandardPublicKey(p2p *P2PService, peerIDStr string) (crypto.PublicKey, error) {
	pubKey, err := GetPeerPublicKey(p2p, peerIDStr)
	if err != nil {
		return nil, err
	}

	// Convert to standard crypto.PublicKey
	stdPubKey, err := libp2pcrypto.PubKeyToStdKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to standard public key: %s", err.Error())
	}

	return stdPubKey, nil
}

// EncryptForPeer encrypts a message for a specific peer using their public key
func EncryptForPeer(p2p *P2PService, message []byte, peerIDStr string) ([]byte, error) {
	stdPubKey, err := GetPeerStandardPublicKey(p2p, peerIDStr)
	if err != nil {
		return nil, err
	}

	// Convert to RSA public key and encrypt
	rsaPubKey, ok := stdPubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("peer's public key is not RSA")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, message)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message: %s", err.Error())
	}

	return ciphertext, nil
}
