package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type HasherInterface interface {
	Hash(str string) (string, error)
	Verify(hashedString string, str string) bool
}

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (*Hasher) Hash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (*Hasher) Verify(hashedString, str string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(str))
	return err == nil
}
