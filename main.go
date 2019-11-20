package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Credentials ...
type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Claims ...
type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.StandardClaims
}

// ErrorResp ...
type ErrorResp struct {
	Message string `json:"message"`
}

var jwtKey = []byte("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")

// TODO: rmeove to Database
var users = map[string]string{
	"test@mail.ru":  "111",
	"user2@mail.ru": "222",
}

// Signin ...
func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{
			Message: err.Error(),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	expectedPassword, ok := users[creds.Email]

	if !ok || expectedPassword != creds.Password {
		resp, _ := json.Marshal(ErrorResp{
			Message: "Incorrect Credentials.",
		})

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{
			Message: err.Error(),
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	resp, _ := json.Marshal(struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	})

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	return
}

func main() {
	http.HandleFunc("/signin", Signin)
	// http.HandleFunc("/welcome", Welcome)
	// http.HandleFunc("/refresh", Refresh)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
