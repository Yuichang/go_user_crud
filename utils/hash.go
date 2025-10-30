package utils

import "golang.org/x/crypto/bcrypt"

// パスワードをbcryptでハッシュ化させる
func GenerateHash(passwd string) (string, error) {
	hashed_passwd, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed_passwd), nil
}

// DBに登録してあるハッシュ値と入力されたパスワードを照合する
func VerifyPassword(inputPasswd string, storedHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(inputPasswd)) == nil
}
