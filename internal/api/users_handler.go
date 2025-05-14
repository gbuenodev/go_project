package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gbuenodev/goProject/internal/store"
	"github.com/gbuenodev/goProject/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
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

func (uh *UserHandler) validateRegisterUserRequest(r *registerUserRequest) error {
	// Username validation
	if len(r.Username) < 3 || len(r.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}
	if !utils.IsValidUsername(r.Username) {
		return errors.New("username can only contain alphanumeric characters and underscores")
	}
	// Check if username already exists
	existingUser, err := uh.userStore.GetUserByUsername(r.Username)
	if err != nil && err != sql.ErrNoRows {
		uh.logger.Printf("Error checking username existence: %v", err)
		return errors.New("internal server error")
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Email validation
	if !utils.IsValidEmail(r.Email) {
		return errors.New("invalid email format")
	}

	// Password validation
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if !utils.IsValidPassword(r.Password) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	// Bio validation
	if len(r.Bio) > 160 {
		return errors.New("bio must be 160 characters or less")
	}

	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = uh.validateRegisterUserRequest(&req)
	if err != nil {
		uh.logger.Printf("ERROR: validating register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// Hash the password
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}
