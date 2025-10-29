package models

import (
	"database/sql"
)

// 引数のユーザー名がDBに存在するか判定する
func (s *Server) UserExistsByName(name string) (bool, error) {
	var user User

	// プレースホルダ用意するとSQL injection対策になる
	row := s.DB.QueryRow("SELECT *FROM users WHERE name = ?", name)

	err := row.Scan(&user.Name)

	// ユーザー名での検索結果が空の場合
	if err == sql.ErrNoRows {
		return false, nil
	} else {
		return true, nil
	}
}
