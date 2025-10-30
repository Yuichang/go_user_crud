package handlers

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/Yuichang/go_user_crud/models"
	"github.com/Yuichang/go_user_crud/utils"
)

var resultTmpl = template.Must(template.ParseFiles("templates/signup_success.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, nil)
}

// サインイン用のハンドラ
func MakeSigninHandler(s *models.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// POSTメソッド以外はエラー処理
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

		hashedPasswd := utils.EasyEncrypt(passwd)

		// ユニークなので、DBにユーザー名とメールを挿入する
		if err := s.InsertUser(name, mail, hashedPasswd); err != nil {
			// メール認証もそうだけど、エラー出た時のhtmlファイルは用意する必要がある
			http.Error(w, "DB登録に失敗しました: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// テンプレートに渡すデータ
		data := struct {
			Name         string
			Mail         string
			HashedPasswd string
		}{
			Name:         name,
			Mail:         mail,
			HashedPasswd: hashedPasswd,
		}

		// 結果のページの表示
		if err := resultTmpl.Execute(w, data); err != nil {
			http.Error(w, "テンプレートの描画に失敗しました", http.StatusInternalServerError)
			return
		}
	}
}

// ログイン用のハンドラ

func MakeLoginHandler(s *models.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTメソッド以外はエラー処理
		if r.Method != http.MethodPost {
			http.Error(w, "POSTメソッドで送信してください", http.StatusMethodNotAllowed)
			return
		}

		// フォームの値の取得
		name := r.FormValue("name")
		hashed_passwd := utils.EasyEncrypt(r.FormValue("passwd"))

		// SQLを使って名前とパスワードが等しいかを判定

		// 名前、ハッシュ化パスワードでDBを検索する
		row := s.DB.QueryRow("SELECT id FROM users WHERE name = ? AND hashed_passwd = ?", name, hashed_passwd)

		var id int
		err := row.Scan(&id)

		// 検索結果が空の場合（名前、パスワードの最低1つは間違っている）
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		} else if err != nil {
			// それ以外のDBエラー
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// ログイン成功なのでページを遷移させる
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

}
