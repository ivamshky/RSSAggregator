package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB  *pgx.Conn
	ctx context.Context
}

func main() {
	godotenv.Load(".env")

	// db
	dbUrl, exists := os.LookupEnv("DB_URL")
	if !exists {
		log.Fatal("DB_URL env not found")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal("couldn't connect err:", err)
	}
	defer conn.Close(context.Background())

	apiCfg := apiConfig{
		DB:  conn,
		ctx: ctx,
	}

	router := configRouter()

	v1Router := chi.NewRouter()
	v1Router.Get("/ready", HandleReadiness)
	v1Router.Get("/err", HandleError)
	v1Router.Post("/create", apiCfg.HandleCreateUser)

	router.Mount("/v1", v1Router)

	// server
	portStr, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("PORT env not found")
	}
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}

	fmt.Println("Server running at port:", portStr)
	srvErr := srv.ListenAndServe()
	if srvErr != nil {
		log.Fatal(err)
	}
}

func configRouter() (r *chi.Mux) {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello"))
	})

	return router
}
