package sessions

import (
	"net/http"
	"time"

	"github.com/GDSC-Phenikaa/ctf-backend/env"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func InitSessionStore() {
	// Initialize the session store with a secret key
	store = sessions.NewCookieStore([]byte(getSessionSecret()))
}

func getSessionSecret() string {
	secret := env.SessionSecret()
	if secret == "" {
		panic("SESSION_SECRET environment variable is not set")
	}
	return secret
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "session")
}

func SetUserID(w http.ResponseWriter, r *http.Request, userID uint) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Values["user_id"] = userID
	return session.Save(r, w)
}

func GetUserID(r *http.Request) (uint, bool) {
	session, err := GetSession(r)
	if err != nil {
		return 0, false
	}
	id, ok := session.Values["user_id"].(uint)
	return id, ok
}

func DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),                     // issued at
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(env.JwtSecret()))
}
