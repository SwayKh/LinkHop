package main

import (
	"log"
	"net/http"
	"os"

	"urlshortner/database"
	"urlshortner/handlers"
	"urlshortner/middleware"
	"urlshortner/render"

	"github.com/joho/godotenv"
)

func auth(next http.HandlerFunc) http.Handler {
	return middleware.AuthMiddleware(next)
}

func main() {
	godotenv.Load()

	database.InitDB()
	render.Load()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /favicon.ico", render.Favicon)
	mux.HandleFunc("GET /{$}", handlers.IndexHandler)
	mux.HandleFunc("GET /signup", handlers.SignupHandler)
	mux.HandleFunc("POST /signup", handlers.SignupHandler)
	mux.HandleFunc("GET /login", handlers.LoginHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)
	mux.HandleFunc("GET /verify", handlers.VerifyHandler)
	mux.HandleFunc("POST /verify", handlers.VerifyHandler)
	mux.HandleFunc("POST /resend-otp", handlers.ResendOTPHandler)

	// Protected routes
	mux.Handle("GET /dashboard", auth(handlers.DashboardHandler))
	mux.Handle("GET /links/new", auth(handlers.CreateLinkHandler))
	mux.Handle("POST /api/links", auth(handlers.CreateLinkHandler))
	mux.Handle("GET /links/{id}/edit", auth(handlers.EditLinkHandler))
	mux.Handle("POST /api/links/{id}/update", auth(handlers.UpdateLinkHandler))
	mux.Handle("POST /api/links/{id}/delete", auth(handlers.DeleteLinkHandler))
	mux.Handle("GET /links/{id}/analytics", auth(handlers.AnalyticsHandler))

	// Short link redirect
	mux.HandleFunc("GET /{shortCode}", handlers.RedirectHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
