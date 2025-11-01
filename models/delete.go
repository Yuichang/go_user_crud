package models

// 引数のユーザーIDのデータをDBから削除する

func (s *Server) UserDeleteByID(ID int) (bool, error) {
	result, err := s.DB.Exec("DELETE FROM users WHERE id=?", ID)
	if err != nil {
		return false, err
	}

	// 削除された行数を確認する
	deleteCnt, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	// 行が削除されなかった場合
	if deleteCnt == 0 {
		return false, nil
	}

	// 削除完了
	return true, nil

}
