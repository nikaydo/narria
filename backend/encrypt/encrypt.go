package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

type Security struct {
	Main       SecurityStruct
	Recovery   SecurityStruct
	MainForPin SecurityStruct
	Pin        SecurityStruct
}

type SecurityStruct struct {
	Salt    string `json:"salt"`
	Wrapped string `json:"wrapped"`
	Nonce   string `json:"nonce"`
	Key     string
}

// InitUser создаёт новый DEK, оборачивает его KEK и возвращает все данные
func InitUser(password string) (Security, []byte, error) {
	dek := make([]byte, 32)
	rand.Read(dek)

	main, err := MakeEncrypt(dek, []byte(password))
	if err != nil {
		return Security{}, nil, err
	}
	recoveryKey := make([]byte, 16)
	rand.Read(recoveryKey)

	recovery, err := MakeEncrypt(dek, []byte(recoveryKey))
	if err != nil {
		return Security{}, nil, err
	}
	return Security{Main: main, Recovery: recovery}, dek, nil
}

// GetDekUser расшифровывает DEK по паролю или recovery key
// For passwords: password is plain text, will be hex-encoded
// For recovery keys: password is already hex-encoded string (32 chars for 16 bytes), use as-is
func GetDekUser(key []byte, security SecurityStruct) ([]byte, error) {
	salt, err := hex.DecodeString(security.Salt)
	if err != nil {
		return nil, err
	}

	kek := deriveKeyFromPassword(key, salt)

	wrapped, err := hex.DecodeString(security.Wrapped)
	if err != nil {
		return nil, err
	}

	nonce, err := hex.DecodeString(security.Nonce)
	if err != nil {
		return nil, err
	}

	dek, err := unwrapDEK(kek, wrapped, nonce)
	if err != nil {
		return nil, err
	}
	return dek, nil
}

// EncryptAES шифрует AES с помощью KEK
func EncryptAES(key, plaintext []byte) ([]byte, error) {
	aesGCM, err := aesGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// DecryptAES расшифровывает AES с помощью KEK
func DecryptAES(key, data []byte) ([]byte, error) {
	aesGCM, err := aesGCM(key)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("data lower than nonce size")
	}

	plaintext, err := aesGCM.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// DeriveKeyFromPassword создаёт KEK (Key Encryption Key)
func deriveKeyFromPassword(key, salt []byte) []byte {
	return argon2.IDKey(key, salt, 1, 64*1024, 4, 32)
}

// WrapDEK шифрует DEK с помощью KEK
func wrapDEK(kek, dek []byte) (wrapped, nonce []byte, err error) {
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	wrapped = aesgcm.Seal(nil, nonce, dek, nil)
	return wrapped, nonce, nil
}

// UnwrapDEK расшифровывает DEK с помощью KEK
func unwrapDEK(kek, wrapped, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != aesgcm.NonceSize() {
		return nil, errors.New("invalid nonce length")
	}

	dek, err := aesgcm.Open(nil, nonce, wrapped, nil)
	if err != nil {
		return nil, err
	}

	return dek, nil
}

func MakeEncrypt(dek []byte, key []byte) (SecurityStruct, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	Kek := deriveKeyFromPassword(key, salt)
	Wrapped, Nonce, err := wrapDEK(Kek, dek)
	if err != nil {
		return SecurityStruct{}, err
	}
	return SecurityStruct{
		Salt:    hex.EncodeToString(salt),
		Wrapped: hex.EncodeToString(Wrapped),
		Nonce:   hex.EncodeToString(Nonce),
		Key:     hex.EncodeToString(key),
	}, nil
}

func aesGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesGCM, nil
}
