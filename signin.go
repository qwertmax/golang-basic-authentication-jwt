package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	return
}
