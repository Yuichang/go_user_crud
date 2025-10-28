package handlers

import (
	"net/http"
	"text/template"

	"github.com/Yuichang/go_user_crud/models"
	"github.com/Yuichang/go_user_crud/utils"
)

var resultTmpl = template.Must(template.ParseFiles("templates/result.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// Serverを受け取ってhttp.HandleFuncを返す
func MakeSubmitHandler(s *models.Server) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// 一旦POSTメソッド以外はエラー処理
		if r.Method != http.MethodPost {
			http.Error(w, "POSTメソッドで送信してください", http.StatusMethodNotAllowed)
			return
		}

		// フォームの値の取得
		name := r.FormValue("name")
		mail := r.FormValue("mail")

		// バリデーションチェック
		if err_mes := utils.CheckValidation(name, mail); err_mes != "OK" {
			http.Error(w, err_mes, http.StatusBadRequest)
			return
		}

		// データベースに名前とメールアドレスを登録する
		// 後で名前がユニークかどうか確認する
		if err := s.InsertUser(name, mail); err != nil {
			http.Error(w, "DB登録に失敗しました: "+err.Error(), http.StatusInternalServerError)
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
}
