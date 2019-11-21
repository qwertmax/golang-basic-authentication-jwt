package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Signup ...
func Signup(w http.ResponseWriter, r *http.Request) {
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

	db, ok := r.Context().Value("db").(*sql.DB)
	if !ok {
		http.Error(w, "could not get database connection pool from context", 500)
		return
	}

	resp, err := db.Query(fmt.Sprintf("SELECT email, password FROM users LIMIT 1"))
	if err != nil {
		fmt.Println(err)
	}
	var email string
	var password string
	resp.Scan(&email, &password)
	fmt.Println(email, password)

	w.WriteHeader(http.StatusOK)
	return
}
