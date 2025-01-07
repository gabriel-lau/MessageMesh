package backend

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"

	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
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
