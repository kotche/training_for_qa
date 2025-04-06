package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

var db *sql.DB

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5433 user=youruser password=yourpassword dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Авто-создание таблицы
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		login TEXT PRIMARY KEY,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	http.HandleFunc("/register", registerHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !isValidPassword(req.Password) {
		http.Error(w, "Пароль слишком простой", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", req.Login, req.Password)
	if err != nil {
		http.Error(w, "Пользователь уже существует в БД", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(Response{Message: "Пользователь успешно зарегистрирован"})
}

func isValidPassword(pw string) bool {
	if len(pw) < 7 {
		return false
	}
	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	special := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return upper.MatchString(pw) && lower.MatchString(pw) && special.MatchString(pw)
}
