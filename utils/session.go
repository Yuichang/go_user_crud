package utils

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const SessionName = "sid"

func NewSessionStore() *sessions.CookieStore {

	hashKey := []byte(os.Getenv("SESS_HASH_KEY"))
	blockKey := []byte(os.Getenv("SESS_BLOCK_KEY"))

	store := sessions.NewCookieStore(hashKey, blockKey)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 24時間
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, //本番はtrueにする
	}
	return store
}

// 未ログインであれば/Loginにリダイレクト

func AuthRequired(store *sessions.CookieStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := store.Get(r, SessionName)
		if _, ok := sess.Values["uid"]; !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
