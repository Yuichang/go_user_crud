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
		passwd := r.FormValue("passwd")

		// バリデーションチェック
		if err_mes := utils.CheckValidation(name, mail); err_mes != "OK" {
			http.Error(w, err_mes, http.StatusBadRequest)
			return
		}

		// 登録する予定のユーザー名がユニークか判定する
		exists, err := s.UserExistsByName(name)

		if err != nil {
			http.Error(w, "ユーザー確認でエラーが発生しました: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// ユーザー名が既に存在している場合はエラーを出す
		if exists {
			http.Error(w, "その名前のユーザーは既に存在しています", http.StatusBadRequest)
			return
		}

		// ユニークなので、DBにユーザー名とメールを挿入する
		if err := s.InsertUser(name, mail); err != nil {
			// メール認証もそうだけど、エラー出た時のhtmlファイルは用意する必要がある
			http.Error(w, "DB登録に失敗しました: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// テンプレートに渡すデータ
		data := struct {
			Name   string
			Mail   string
			Passwd string
		}{
			Name:   name,
			Mail:   mail,
			Passwd: passwd,
		}

		// 結果のページの表示
		if err := resultTmpl.Execute(w, data); err != nil {
			http.Error(w, "テンプレートの描画に失敗しました", http.StatusInternalServerError)
			return
		}

	}
}
