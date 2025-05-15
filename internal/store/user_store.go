package store

import (
	"database/sql"
	_ "golang.org/x/crypto/bcrypt"
	"time"
)

type password struct {
	plaintext *string
	hash      []byte
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"` // Don't include password hash in JSON response
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, avatar_url, bio)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash, user.AvatarURL, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT id, username, email, password_hash, avatar_url, bio, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.AvatarURL,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET username = $1, email = $2, avatar_url = $3, bio = $4, updated_at = NOW()
	WHERE id = $4
	RETURNING updated_at
	`
	result, err := s.db.Exec(query, user.Username, user.Email, user.AvatarURL, user.Bio, user.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
