package models

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Server struct {
	DB *sql.DB
}

func Connect() (*sql.DB, error) {
	// 環境変数に.envファイルの内容を登録する
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error open .env file")
		return nil, err
	}

	// 環境変数に登録されている値を入れる
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_ADDRESS"),
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("DB open Error")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		fmt.Println("Ping error", err)
		return nil, err
	}
	fmt.Println("connection success!!")
	return db, nil
}
