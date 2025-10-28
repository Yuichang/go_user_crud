package main

import (
	"log"
	"net/http"

	"github.com/Yuichang/go_user_crud/handlers"
	"github.com/Yuichang/go_user_crud/models"
)

func main() {
	db, err := models.Connect()
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	defer db.Close()
	log.Println("Connection Success!")

	s := models.NewServer(db)

	// 静的ファイルのハンドラ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 各ページのハンドラ
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/submit", handlers.MakeSubmitHandler(s))

	log.Println("Server Start...: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
