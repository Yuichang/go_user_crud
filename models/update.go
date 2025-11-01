package models

func (s *Server) UpdateNameByID(ID int, newName string) (bool, error) {
	res, err := s.DB.Exec("UPDATE users SET name = ? WHERE id = ?", newName, ID)

	if err != nil {
		return false, err
	}
	changeCnt, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return changeCnt > 0, nil
}
