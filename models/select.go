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

// 引数のユーザーIDのユーザー名を返す
func (s *Server) ReturnNameByID(ID int) (string, error) {
	var name string
	row := s.DB.QueryRow("SELECT name FROM users WHERE ID = ?", ID)
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
