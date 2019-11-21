package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Signin ...
func Signin(w http.ResponseWriter, r *http.Request) {
	var userAuth User
	err := json.NewDecoder(r.Body).Decode(&userAuth)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{
			Message: err.Error(),
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	db, ok := r.Context().Value("db").(*sql.DB)
	if !ok {
		http.Error(w, "could not get database connection pool from context", 500)
		return
	}

	var user User
	err = user.Get(db, userAuth.Email)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if !user.IsHashEqual(userAuth.CreateHash(userAuth.Password)) {
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
		Email: user.Email,
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	return
}
