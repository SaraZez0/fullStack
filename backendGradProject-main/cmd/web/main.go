package main

import (
	"database/sql"
	"gp-backend/database"
	"log"
	"net/http"
	"time"
	// "github.com/alexedwards/scs/postgresstore"
	// "github.com/alexedwards/scs/v2"
	_ "github.com/lib/pq"
	"github.com/r3labs/sse/v2"
)
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // السماح فقط لـ frontend
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // في حال كنت تستخدم cookies أو authentication
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
type application struct {
	sse *sse.Server
	db *sql.DB
	udb database.UserDB
	// sessionManager *scs.SessionManager
}

func main() {
	app, err := NewApplication("")
	if err != nil {
		log.Fatalln(err.Error())
	}

	server := &http.Server{
		Addr:         ":5000",
		Handler:      corsMiddleware(app.routes()), // تمرير Middleware CORS
		ReadTimeout:  10 * time.Second,             // مهلة القراءة
		WriteTimeout: 10 * time.Second,             // مهلة الكتابة
		IdleTimeout:  30 * time.Second,             // مهلة الاتصال غير النشط
	}

	log.Println("SERVER STARTED")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}


func OpenDB() (*sql.DB, error) {
	// TODO: connect through ssl (enable sslmode)
	db, err := sql.Open("postgres", "host=localhost user=emergency password=emergency dbname=cars sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	log.Println("CONNECTED TO DATABASE CARS")

	return db, nil
}


func NewApplication(sseParameter string) (*application, error) {
	if sseParameter == "" {
		sseParameter = "messages"
	}

	// Creating the SSE server
	server := sse.New()
	server.CreateStream(sseParameter)

	// db, err := OpenDB()
	// if err != nil {
	// 	return nil, err
	// }
	// udb := database.UserDB{
	// 	DB: db,
	// }

	// sessionManager := scs.New()
	// sessionManager.Store = postgresstore.New(db)
	// sessionManager.Lifetime = 12 * time.Hour
	return &application{
		sse: server,
		// db: db,
		// udb: udb,
		// sessionManager: sessionManager,
	}, nil
}
