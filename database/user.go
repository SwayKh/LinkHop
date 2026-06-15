package database

import "time"

type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	IsVerified   bool
	CreatedAt    string
}

type VerificationCode struct {
	ID        int64
	Email     string
	Code      string
	ExpiresAt time.Time
	Used      bool
}

func CreateUser(email, username, passwordHash string) error {
	_, err := DB.Exec(
		"INSERT INTO users (email, username, password_hash, is_verified) VALUES ($1, $2, $3, false)",
		email, username, passwordHash,
	)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, email, username, password_hash, is_verified, created_at FROM users WHERE email = $1", email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsVerified, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, email, username, password_hash, is_verified, created_at FROM users WHERE username = $1", username,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsVerified, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, email, username, password_hash, is_verified, created_at FROM users WHERE id = $1", id,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.IsVerified, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func SetUserVerified(email string) error {
	_, err := DB.Exec("UPDATE users SET is_verified = true WHERE email = $1", email)
	return err
}

func CreateVerificationCode(email, code string, expiresAt time.Time) error {
	_, err := DB.Exec(
		"INSERT INTO verification_codes (email, code, expires_at) VALUES ($1, $2, $3)",
		email, code, expiresAt,
	)
	return err
}

func GetValidVerificationCode(email, code string) (*VerificationCode, error) {
	vc := &VerificationCode{}
	err := DB.QueryRow(
		"SELECT id, email, code, expires_at, used FROM verification_codes WHERE email = $1 AND code = $2 AND used = false AND expires_at > NOW() ORDER BY id DESC LIMIT 1",
		email, code,
	).Scan(&vc.ID, &vc.Email, &vc.Code, &vc.ExpiresAt, &vc.Used)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

func MarkVerificationCodeUsed(id int64) error {
	_, err := DB.Exec("UPDATE verification_codes SET used = true WHERE id = $1", id)
	return err
}
