// backend/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Логин, пароль и имя БД приходят ТОЛЬКО из окружения (файл .env),
	// значений по умолчанию для секретов в коде нет
	user := mustEnv("DB_USER")
	password := mustEnv("DB_PASSWORD")
	dbname := mustEnv("DB_NAME")

	dsn := "host=" + env("DB_HOST", "db") +
		" port=" + env("DB_PORT", "5432") +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbname +
		" sslmode=disable"

	// sql.Open соединение не устанавливает, только готовит пул
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	staticDir := env("STATIC_DIR", "/app/static")

	// /health: приложение живо И база отвечает (используется в тестах CI)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "db unavailable: "+err.Error(), http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "ok")
	})

	// Все записи — их показывает веб-страница
	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.QueryContext(r.Context(), `SELECT id, name FROM items ORDER BY id`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := []Item{}
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, it)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	// Количество записей (для сравнения с БД в тестах CI)
	http.HandleFunc("/api/items/count", func(w http.ResponseWriter, r *http.Request) {
		var n int
		if err := db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM items`).Scan(&n); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, n)
	})

	// FileServer сам отдаёт index.html для "/"
	http.Handle("/", http.FileServer(http.Dir(staticDir)))

	addr := ":" + env("HTTP_PORT", "8080")
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// mustEnv завершает программу, если обязательная переменная не задана
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set (see .env)", key)
	}
	return v
}
