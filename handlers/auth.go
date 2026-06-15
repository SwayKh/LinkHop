package handlers

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"urlshortner/database"
	"urlshortner/middleware"
	"urlshortner/render"

	"golang.org/x/crypto/bcrypt"

	emailpkg "urlshortner/email"
)

func generateOTP() (string, error) {
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code[i] = byte('0') + byte(n.Int64())
	}
	return string(code), nil
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		render.Template(w, "signup", &render.Data{})
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		render.Template(w, "signup", &render.Data{Error: "All fields are required"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		render.Template(w, "signup", &render.Data{Error: "Failed to create account"})
		return
	}

	if err := database.CreateUser(email, username, string(hash)); err != nil {
		render.Template(w, "signup", &render.Data{Error: "Email or username already taken"})
		return
	}

	otp, err := generateOTP()
	if err != nil {
		render.Template(w, "signup", &render.Data{Error: "Failed to generate verification code"})
		return
	}

	if err := database.CreateVerificationCode(email, otp, time.Now().Add(15*time.Minute)); err != nil {
		render.Template(w, "signup", &render.Data{Error: "Failed to create verification code"})
		return
	}

	if err := emailpkg.SendOTP(email, otp); err != nil {
		log.Printf("Failed to send OTP to %s: %v", email, err)
		render.Template(w, "signup", &render.Data{Error: "Failed to send verification email"})
		return
	}

	http.Redirect(w, r, "/verify?email="+email, http.StatusSeeOther)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	if r.Method == "GET" {
		if email == "" {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		render.Template(w, "verify", &render.Data{FormEmail: email})
		return
	}

	code := strings.TrimSpace(r.FormValue("code"))
	email = strings.TrimSpace(r.FormValue("email"))

	if email == "" || code == "" {
		render.Template(w, "verify", &render.Data{FormEmail: email, Error: "Code is required"})
		return
	}

	vc, err := database.GetValidVerificationCode(email, code)
	if err != nil {
		render.Template(w, "verify", &render.Data{FormEmail: email, Error: "Invalid or expired code"})
		return
	}

	if err := database.MarkVerificationCodeUsed(vc.ID); err != nil {
		render.Template(w, "verify", &render.Data{FormEmail: email, Error: "Failed to verify code"})
		return
	}

	if err := database.SetUserVerified(email); err != nil {
		render.Template(w, "verify", &render.Data{FormEmail: email, Error: "Failed to verify account"})
		return
	}

	http.Redirect(w, r, "/login?success=Email+verified+successfully", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		success := r.URL.Query().Get("success")
		render.Template(w, "login", &render.Data{Success: success})
		return
	}

	login := strings.TrimSpace(r.FormValue("login"))
	password := r.FormValue("password")

	if login == "" || password == "" {
		render.Template(w, "login", &render.Data{Error: "Email or username and password are required"})
		return
	}

	var user *database.User
	var err error

	if strings.Contains(login, "@") {
		user, err = database.GetUserByEmail(login)
	} else {
		user, err = database.GetUserByUsername(login)
	}

	if err != nil {
		render.Template(w, "login", &render.Data{Error: "Invalid email/username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		render.Template(w, "login", &render.Data{Error: "Invalid email/username or password"})
		return
	}

	if !user.IsVerified {
		http.Redirect(w, r, "/verify?email="+user.Email, http.StatusSeeOther)
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		render.Template(w, "login", &render.Data{Error: "Failed to create session"})
		return
	}

	secure := os.Getenv("COOKIE_SECURE") == "true"

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	secure := os.Getenv("COOKIE_SECURE") == "true"

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ResendOTPHandler(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.FormValue("email"))
	if email == "" {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	otp, err := generateOTP()
	if err != nil {
		http.Redirect(w, r, "/verify?email="+email+"&error=Failed+to+generate+code", http.StatusSeeOther)
		return
	}

	if err := database.CreateVerificationCode(email, otp, time.Now().Add(15*time.Minute)); err != nil {
		http.Redirect(w, r, "/verify?email="+email+"&error=Failed+to+create+code", http.StatusSeeOther)
		return
	}

	if err := emailpkg.SendOTP(email, otp); err != nil {
		http.Redirect(w, r, "/verify?email="+email+"&error=Failed+to+send+email", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/verify?email="+email+"&success=Code+resent", http.StatusSeeOther)
}
