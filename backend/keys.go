package backend

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateKeyFiles() error {
	_ = os.Mkdir("./keys", os.ModePerm)
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key pair:", err)
		return err
	}

	// Save the private key to a file
	privateKeyFile, err := os.Create("keys/private_key.pem")
	if err != nil {
		fmt.Println("Error creating private key file:", err)
		return err
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		fmt.Println("Error encoding private key to file:", err)
		return err
	}

	// Save the public key to a file
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Error marshaling public key:", err)
		return err
	}

	publicKeyFile, err := os.Create("keys/public_key.pem")
	if err != nil {
		fmt.Println("Error creating public key file:", err)
		return err
	}
	defer publicKeyFile.Close()

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		fmt.Println("Error encoding public key to file:", err)
		return err
	}

	fmt.Println("Private and public keys generated and saved to private_key.pem and public_key.pem")
	return nil
}

func ReadPrivateKeyFile() (*rsa.PrivateKey, error) {
	// Read the private key from file
	privateKeyFile, err := os.Open("keys/private_key.pem")
	if err != nil {
		fmt.Println("Error opening private key file:", err)
		return nil, err
	}
	defer privateKeyFile.Close()

	privateKeyPEM, err := os.ReadFile("keys/private_key.pem")
	if err != nil {
		fmt.Println("Error reading private key file:", err)
		return nil, err
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		fmt.Println("Error decoding private key PEM block")
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return nil, err
	}

	return privateKey, nil
}

func ReadPublicKeyFile() (*rsa.PublicKey, error) {
	// Read the public key from file
	publicKeyFile, err := os.Open("keys/public_key.pem")
	if err != nil {
		fmt.Println("Error opening public key file:", err)
		return nil, err
	}
	defer publicKeyFile.Close()

	publicKeyPEM, err := os.ReadFile("keys/public_key.pem")
	if err != nil {
		fmt.Println("Error reading public key file:", err)
		return nil, err
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		fmt.Println("Error decoding public key PEM block")
		return nil, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key:", err)
		return nil, err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Error converting public key to RSA public key")
		return nil, err
	}

	return rsaPublicKey, nil
}

func EncryptWithPublicKey(pubKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, message)
}

func DecryptWithPrivateKey(privKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, ciphertext)
}

func SignWithPrivateKey(privKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, privKey, 0, message)
}

func VerifyWithPublicKey(pubKey *rsa.PublicKey, message, signature []byte) error {
	return rsa.VerifyPKCS1v15(pubKey, 0, message, signature)
}
