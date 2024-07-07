package common

import "golang.org/x/crypto/bcrypt"

// BcryptHash 使用bcrypt加密密码
func BcrypaHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes)
}

// BcryptCheck 对比铭文密码和数据库的哈希值 使用bcrypt加密密码
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("password"))
	return err == nil
}
