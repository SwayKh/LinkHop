package database

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    string
}

func CreateUser(email, passwordHash string) error {
	_, err := DB.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", email, passwordHash)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, email, password_hash, created_at FROM users WHERE email = $1", email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		"SELECT id, email, password_hash, created_at FROM users WHERE id = $1", id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
