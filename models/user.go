package models

import "time"

type User struct {
	ID           int64     // db->id
	Name         string    // db->name
	Email        string    // db->email
	HashedPasswd string    // db->hashed_passwd
	CreatedAt    time.Time // db->created_at
	UpdateAt     time.Time // db->updated_at
}
