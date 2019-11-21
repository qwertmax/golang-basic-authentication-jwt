package main

import (
	"context"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Credentials ...
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

// ContextInjector ...
type ContextInjector struct {
	ctx context.Context
	h   http.Handler
}

func (ci *ContextInjector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ci.h.ServeHTTP(w, r.WithContext(ci.ctx))
}
func main() {
	db := Storage{}
	err := db.Connect(GetDBCredentials("config.yml"))
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.DB.Close()

	db.Init()

	ctx := context.WithValue(context.Background(), "db", db.DB)

	http.Handle("/signin", &ContextInjector{ctx, http.HandlerFunc(Signin)})
	http.Handle("/signup", &ContextInjector{ctx, http.HandlerFunc(Signup)})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
