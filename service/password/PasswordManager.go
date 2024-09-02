package password

import (
	"crypto/sha256"
	"encoding/hex"
)

type PasswordManager interface {
	VerifyPassword(password string, expect string) bool
	HashPassword(password string) string
}

type Sha256PasswordManager struct {
	secret string
}

func NewSha256Hash(secret string) PasswordManager {
	return &Sha256PasswordManager{
		secret: secret,
	}
}

func (manager *Sha256PasswordManager) VerifyPassword(password string, expect string) bool {
	actual := manager.HashPassword(password)
	return actual == expect
}

func (manager *Sha256PasswordManager) HashPassword(password string) string {
	password = password + manager.secret
	hashed := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashed[:])
}
