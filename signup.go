package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Signup ...
func Signup(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
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

	err = user.Save(db)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{
			Message: "Email already exist",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	return
}
