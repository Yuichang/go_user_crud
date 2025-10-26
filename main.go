package main

import (
	"log"
	"net/http"
	"regexp"
	"text/template"

	"github.com/Yuichang/go_user_crud/models"
)

var resultTmpl = template.Must(template.ParseFiles("templates/result.html"))

func main() {
	db, err := models.Connect()
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	defer db.Close()
	log.Println("Connection Success!")

	// 静的ファイルのハンドラ
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 各ページのハンドラ
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", submitHandler)

	log.Println("Server Start...: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {

	// 一旦POSTメソッド以外はエラー処理
	if r.Method != http.MethodPost {
		http.Error(w, "POSTメソッドで送信してください", http.StatusMethodNotAllowed)
		return
	}

	// フォームの値の取得
	name := r.FormValue("name")
	mail := r.FormValue("mail")

	// バリデーションチェック
	if err_mes := CheckValidation(name, mail); err_mes != "OK" {
		http.Error(w, err_mes, http.StatusBadRequest)
		return
	}

	// テンプレートに渡すデータ
	data := struct {
		Name string
		Mail string
	}{
		Name: name,
		Mail: mail,
	}

	// 結果のページの表示

	if err := resultTmpl.Execute(w, data); err != nil {
		http.Error(w, "テンプレートの描画に失敗しました", http.StatusInternalServerError)
		return
	}
}

// 名前、メールのバリデーションチェック
func CheckValidation(name string, mail string) string {

	if name == "" {
		return "ユーザーネームを入力してください"
	} else if mail == "" {
		return "メールアドレスを入力してください"
	} else if !MailValidation(mail) {
		return "メールアドレスの形式が間違っています"
	} else {
		return "OK"
	}
}

// メールのバリデーションチェック
func MailValidation(mail string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(mail)
}
