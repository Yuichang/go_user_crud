package models

import "fmt"

func (s *Server) InsertUser(name string, mail string, storedHash string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("transaction start err:%v", err)
	}

	defer tx.Rollback()

	// 後でユーザーネームがユニークかチェックする
	if _, err = tx.Exec("INSERT INTO users (name,mail,stored_hash)VALUES(?,?,?)", name, mail, storedHash); err != nil {
		return fmt.Errorf("insert user error: %v", err)
	}

	return tx.Commit()

}
