package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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
	http.HandleFunc("/login", loginHandler)

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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	var password string
	err := db.QueryRow("SELECT password FROM users WHERE login=$1", req.Login).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		} else {
			http.Error(w, "Ошибка поиска пользователя", http.StatusInternalServerError)
		}
		return
	}

	// bug
	if req.Password != "" && password != req.Password {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(Response{Message: "Успешный вход"})
}

func isValidPassword(pw string) bool {
	// bug
	if len(pw) < 7 {
		return false
	}
	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	special := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return upper.MatchString(pw) && lower.MatchString(pw) && special.MatchString(pw)
}
