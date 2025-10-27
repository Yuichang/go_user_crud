package utils

import "regexp"

// 名前、メールのバリデーションチェック
func CheckValidation(name string, mail string) string {

	if name == "" {
		return "ユーザーネームを入力してください"
	} else if mail == "" {
		return "メールアドレスを入力してください"
	} else if !MailValidation(mail) {
		return "メールアドレスの形式が間違っています"
	} else {
		return "OK"
	}
}

// メールのバリデーションチェック
func MailValidation(mail string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(mail)
}
