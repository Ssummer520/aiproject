package infrastructure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"
)

const piiKeyEnv = "USER_PII_ENCRYPTION_KEY"

type PIICrypto struct {
	key []byte
}

func NewPIICrypto() (*PIICrypto, error) {
	raw := strings.TrimSpace(os.Getenv(piiKeyEnv))
	if raw != "" {
		key, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return nil, err
		}
		if len(key) != 32 {
			return nil, errors.New("USER_PII_ENCRYPTION_KEY must be a base64 encoded 32-byte key")
		}
		return &PIICrypto{key: key}, nil
	}
	sum := sha256.Sum256([]byte("travel-api-demo-pii-key"))
	return &PIICrypto{key: sum[:]}, nil
}

func (c *PIICrypto) Encrypt(plain string) (string, error) {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return "", nil
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

func (c *PIICrypto) Hash(documentType, documentNo string) string {
	message := strings.ToUpper(strings.TrimSpace(documentType)) + ":" + strings.ToUpper(strings.TrimSpace(documentNo))
	mac := hmac.New(sha256.New, c.key)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func MaskDocument(documentNo string) string {
	documentNo = strings.TrimSpace(documentNo)
	if documentNo == "" {
		return ""
	}
	runes := []rune(documentNo)
	if len(runes) <= 4 {
		return strings.Repeat("*", len(runes))
	}
	if len(runes) <= 8 {
		return string(runes[:1]) + strings.Repeat("*", len(runes)-2) + string(runes[len(runes)-1:])
	}
	return string(runes[:2]) + strings.Repeat("*", len(runes)-6) + string(runes[len(runes)-4:])
}
