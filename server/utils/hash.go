package utils

import "golang.org/x/crypto/bcrypt"

func BcryptHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
