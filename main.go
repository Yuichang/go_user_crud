package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	// 静的ファイルのハンドラ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 各ページのハンドラ
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", submitHnadler)

	log.Println("Server Start...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func submitHnadler(w http.ResponseWriter, r *http.Request) {

	// 一旦POSTメソッド以外はエラー処理
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドで送信してください", http.StatusMethodNotAllowed)
		return
	}

	// フォームの値の取得
	name := r.FormValue("name")
	mail := r.FormValue("mail")

	// バリデーション
	
}
