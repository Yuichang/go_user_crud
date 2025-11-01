package main

import (
	"log"
	"net/http"

	"github.com/Yuichang/go_user_crud/handlers"
	"github.com/Yuichang/go_user_crud/models"
	"github.com/Yuichang/go_user_crud/utils"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db, err := models.Connect()
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	defer db.Close()
	log.Println("Connection Success!")

	s := models.NewServer(db)
	store := utils.NewSessionStore()
	handlers.SetSessionStore(store)

	// 静的ファイルのハンドラ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 各ページのハンドラ
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	// 認証系
	http.HandleFunc("/login/submit", handlers.MakeLoginHandler(s))
	http.HandleFunc("/logout", handlers.MakeLogoutHandler(store))
	http.HandleFunc("/register/submit", handlers.MakeRegisterHandler(s))

	// 認証必須ページ
	http.Handle("/mypage", utils.AuthRequired(store, http.HandlerFunc(handlers.MypageHandler)))

	log.Println("Server Start...: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}


