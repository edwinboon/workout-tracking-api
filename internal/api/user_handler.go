package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/edwinboon/workout-tracking-api/internal/store"
	"github.com/edwinboon/workout-tracking-api/internal/utils"
)

type RegisterUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
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

	// @TODO - Add password complexity requirements
	if len(req.Password) < 8 || len(req.Password) > 20 {
		return errors.New("password must be between 8 and 20 characters long")
	}

	avatarURLRegex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)

	if req.AvatarURL != "" && !avatarURLRegex.MatchString(req.AvatarURL) {
		return errors.New("invalid avatar URL format")
	}

	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		uh.logger.Printf("ERROR: decodingRegisterUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = uh.ValidateRegisterRequest(&req)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// Hash the password
	err = user.PasswordHash.SetPassword(req.Password)

	if err != nil {
		uh.logger.Printf("ERROR: hashingPassword: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = uh.userStore.CreateUser(user)

	if err != nil {
		uh.logger.Printf("ERROR: registering user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})

}
