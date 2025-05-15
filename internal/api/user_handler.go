package api

import (
	"errors"
	"log"
	"regexp"

	"github.com/edwinboon/workout-tracking-api/internal/store"
)

type RegisterUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	avatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) ValidateRegisterRequest(req *RegisterUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) < 3 || len(req.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters long")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	passwordRegex := regexp.MustCompile(`^(?=.*[A-Z])(?=.*[0-9])(?=.*[^A-Za-z0-9])(?=.{8,}).*$`)

	if !passwordRegex.MatchString(req.Password) {
		return errors.New("password must be at least 8 characters long, contain at least one uppercase letter, one number, and one special character")
	}

	avatarURLRegex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)

	if req.avatarURL != "" && !avatarURLRegex.MatchString(req.avatarURL) {
		return errors.New("invalid avatar URL format")
	}

	return nil

}
