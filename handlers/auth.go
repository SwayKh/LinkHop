package handlers

import (
	"net/http"
	"strings"

	"urlshortner/database"
	"urlshortner/middleware"
	"urlshortner/render"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render.Template(w, "signup", &render.Data{})
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || password == "" {
		render.Template(w, "signup", &render.Data{Error: "Email and password are required"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		render.Template(w, "signup", &render.Data{Error: "Failed to create account"})
		return
	}

	if err := database.CreateUser(email, string(hash)); err != nil {
		render.Template(w, "signup", &render.Data{Error: "Email already exists"})
		return
	}

	http.Redirect(w, r, "/login?success=Account+created+successfully", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		success := r.URL.Query().Get("success")
		render.Template(w, "login", &render.Data{Success: success})
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if email == "" || password == "" {
		render.Template(w, "login", &render.Data{Error: "Email and password are required"})
		return
	}

	user, err := database.GetUserByEmail(email)
	if err != nil {
		render.Template(w, "login", &render.Data{Error: "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		render.Template(w, "login", &render.Data{Error: "Invalid email or password"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		render.Template(w, "login", &render.Data{Error: "Failed to create session"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   86400,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

