package services

import (
	"errors"
	"strings"

	"ticket-system/models"
	"ticket-system/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{db: db, jwtSecret: jwtSecret}
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	if err := utils.ValidateName(input.Name); err != nil {
		return nil, err
	}
	if err := utils.ValidateEmail(input.Email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassword(input.Password); err != nil {
		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	var existing models.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("Email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         strings.TrimSpace(input.Name),
		Email:        email,
		PasswordHash: hash,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(input LoginInput) (string, error) {
	if err := utils.ValidateEmail(input.Email); err != nil {
		return "", err
	}
	if err := utils.ValidatePassword(input.Password); err != nil {
		return "", err
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("Invalid email or password")
		}
		return "", err
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return "", errors.New("Invalid email or password")
	}

	return utils.GenerateToken(user.ID, s.jwtSecret)
}
