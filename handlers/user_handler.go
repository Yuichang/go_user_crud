package handlers

import (
	"net/http"
	"text/template"

	"github.com/Yuichang/go_user_crud/models"
	"github.com/Yuichang/go_user_crud/utils"
	"github.com/gorilla/sessions"
)

var sessionStore *sessions.CookieStore

func SetSessionStore(s *sessions.CookieStore) { sessionStore = s }

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
func MakeRegisterHandler(s *models.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// POSTメソッド以外はエラー処理
		if r.Method != http.MethodPost {
			http.Error(w, "Please use POST method", http.StatusMethodNotAllowed)
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
			http.Error(w, "Error occurred checking existing user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// ユーザー名が既に存在している場合はエラーを出す
		if exists {
			http.Error(w, "This username already exists", http.StatusBadRequest)
			return
		}

		storedHash, err := utils.GenerateHash(passwd)

		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// ユニークなので、DBにユーザー名とメールを挿入する
		if err := s.InsertUser(name, mail, storedHash); err != nil {

			// メール認証もそうだけど、エラー出た時のhtmlファイルは用意する必要がある
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// テンプレートに渡すデータ
		data := struct {
			Name       string
			Mail       string
			StoredHash string
		}{
			Name:       name,
			Mail:       mail,
			StoredHash: storedHash,
		}

		// 結果のページの表示
		if err := resultTmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}

// ログイン用のハンドラ

func MakeLoginHandler(s *models.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POSTメソッド以外はリダイレクト
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// フォームの値の取得
		name := r.FormValue("name")
		passwd := r.FormValue("passwd")

		// 名前でDB検索してid,ハッシュ済みパスワードを持ってくる
		row := s.DB.QueryRow("SELECT id,stored_hash FROM users WHERE name = ?", name)

		var id int
		var storedHash string
		err := row.Scan(&id, &storedHash)

		// DBエラー
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if !utils.VerifyPassword(passwd, storedHash) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		if sessionStore != nil {
			sess, _ := sessionStore.Get(r, utils.SessionName)
			sess.Values["uid"] = id
			if err := sess.Save(r, w); err != nil {
				http.Error(w, "failed to save sesson", http.StatusInternalServerError)
				return
			}
		}

		// 今後セッション認証にする
		http.Redirect(w, r, "/mypage", http.StatusSeeOther)
	}
}

func MakeLogoutHandler(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := store.Get(r, utils.SessionName)
		sess.Options.MaxAge = -1 //セッション削除
		_ = sess.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}

func MakeMypageHandler(s *models.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if sessionStore == nil {
			http.Error(w, "session not ready", http.StatusInternalServerError)
			return
		}

		sess, _ := sessionStore.Get(r, utils.SessionName)
		uid, ok := sess.Values["uid"].(int)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		name, err := s.ReturnNameByID(uid)
		if err != nil {
			http.Error(w, "failed to fetch user name", http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/mypage.html"))
		_ = tmpl.Execute(w, map[string]any{
			"UserName": name,
		})

	}
}
