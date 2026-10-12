package services

import (
	"errors"
	"simple-auction-system/internal/models"
	"simple-auction-system/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key_change_this_in_production")

// AuthService định nghĩa logic nghiệp vụ
type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type authService struct {
	repo repositories.UserRepository
}

// NewAuthService khởi tạo service
func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(username, password string) error {
	// 1. Kiểm tra user đã tồn tại chưa
	_, err := s.repo.GetUserByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}

	// 2. Băm mật khẩu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. Lưu vào DB
	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return s.repo.CreateUser(user)
}

func (s *authService) Login(username, password string) (string, error) {
	// 1. Tìm user trong DB
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. So sánh mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Tạo JWT Token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token hết hạn sau 24h
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString(jwtKey)
}