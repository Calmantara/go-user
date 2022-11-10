package crypto

import "golang.org/x/crypto/bcrypt"

type Hash interface {
	GeneratePassword(password string) (string, error)
	ComparePassword(password string, hashPassword string) error
}

type HashImpl struct{}

func NewHash() Hash {
	return &HashImpl{}
}

func (h *HashImpl) GeneratePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}
func (h *HashImpl) ComparePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
