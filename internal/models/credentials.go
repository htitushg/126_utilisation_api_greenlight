package models

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

var pepper = []byte{0x7, 0x4e, 0x6, 0xe8, 0xc5, 0xc4, 0xf3, 0xe, 0x6d, 0x4d, 0xe, 0x72, 0x8a, 0xbd, 0x85, 0x9c, 0xc8, 0xa9, 0xc7, 0xe3, 0x59, 0x4f, 0x97, 0xe2, 0xb, 0x85, 0x3e, 0x21, 0xad, 0xba, 0xe2, 0x17, 0x13, 0xb8, 0x3f, 0xd1, 0x52, 0x50, 0x6e, 0xa8, 0xd2, 0x8, 0xd3, 0x8a, 0x7f, 0x28, 0xc5, 0xc2, 0x3f, 0x64, 0x99, 0xb8, 0x23, 0x66, 0x11, 0xf0, 0xc4, 0x4, 0x59, 0x3, 0x2d, 0x45, 0x4e, 0xe7}

func newSalt() []byte {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}

func hash(pwd string, salt []byte) string {
	saltedPwd := append(salt, []byte(pwd)...)
	sha := sha512.New()
	sha.Write(append(saltedPwd, pepper...))
	// key stretching
	for i := 0; i <= 10071; i++ {
		// Effectuez une opération sur chaque nombre
		sha.Sum(sha.Sum(nil))
	}
	return hex.EncodeToString(sha.Sum(nil))
}

func NewPwd(pwd string) (string, string) {
	salt := newSalt()
	return hash(pwd, salt), base64.StdEncoding.EncodeToString(salt)
}

func (m *UsersModel) CheckPwd(cred Credentials) bool {
	user, ok := m.GetUser(cred.Username)
	if !ok {
		return false
	}
	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return false
	}
	return user.HashedPwd == hash(cred.Password, salt)
}
