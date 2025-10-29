package utils

import (
	"encoding/base64"
	"fmt"
)

// 適当なkeyを作ってxorしてく
func XorTransform(passwd, key string) string {
	res := []byte(passwd)
	for i := 0; i < len(passwd); i++ {
		res[i] ^= key[i%len(key)]
	}
	return string(res)
}

// 実務ではもちろん使わない
func EasyEncrypt(passwd string) string {
	key := "FZJESOFDSK"
	enc := XorTransform(passwd, key)

	// base64に変換する
	enc64 := base64.StdEncoding.EncodeToString([]byte(enc))

	// 暗号化できてるかのテスト
	fmt.Println("暗号化結果: ", enc64)
	return enc64
}
