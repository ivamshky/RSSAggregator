package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB  *pgx.Conn
	ctx context.Context
}

const DB_URL = "postgresql://postgres:%s@%s:%s/%s?sslmode=disable"

func main() {
	godotenv.Load(".env")

	// db
	dbUrl := fmt.Sprintf(DB_URL, os.Getenv("DB_PWD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_SCHEMA"))
	fmt.Println("Connecting to db", dbUrl)

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("couldn't connect err:", err)
	}
	defer conn.Close(context.Background())

	apiCfg := apiConfig{
		DB:  conn,
		ctx: context.Background(),
	}

	router := configRouter()

	v1Router := chi.NewRouter()
	v1Router.Get("/ready", HandleReadiness)
	v1Router.Get("/err", HandleError)
	v1Router.Post("/create", apiCfg.HandleCreateUser)
	v1Router.Get("/user/{name}", apiCfg.HandleGetByName)

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
