package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// Encryptor handles encryption and decryption operations
type Encryptor struct {
	key []byte
}

// NewEncryptor creates a new encryptor with the given key
// The key will be hashed to ensure it's 32 bytes (256 bits)
func NewEncryptor(key string) *Encryptor {
	hash := sha256.Sum256([]byte(key))
	return &Encryptor{
		key: hash[:],
	}
}

// Encrypt encrypts plaintext using AES-256-GCM
// Returns base64 encoded ciphertext
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create a nonce (number used once)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for safe storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64 encoded ciphertext using AES-256-GCM
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and encrypted data
	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// EncryptIfEnabled encrypts data only if encryptor is not nil
func EncryptIfEnabled(encryptor *Encryptor, plaintext string) (string, error) {
	if encryptor == nil {
		return plaintext, nil
	}
	return encryptor.Encrypt(plaintext)
}

// DecryptIfEnabled decrypts data only if encryptor is not nil
// If data is not valid base64 (old unencrypted data), returns it as-is
func DecryptIfEnabled(encryptor *Encryptor, ciphertext string) (string, error) {
	if encryptor == nil {
		return ciphertext, nil
	}

	// Check if data looks like base64 encoded encrypted data
	// If it's not valid base64, it's probably old unencrypted data
	_, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		// Not valid base64, assume it's old unencrypted data
		return ciphertext, nil
	}

	// Try to decrypt - if it fails, return original (might be unencrypted)
	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		// Decryption failed, might be old unencrypted data
		// Return original without error to maintain backward compatibility
		return ciphertext, nil
	}

	return decrypted, nil
}
