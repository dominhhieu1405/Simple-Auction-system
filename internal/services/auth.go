package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"simple-auction-system/internal/models"
	"simple-auction-system/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidToken       = errors.New("invalid token")
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthService định nghĩa logic nghiệp vụ
type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	ParseToken(token string) (*Claims, error)
}

type authService struct {
	repo   repositories.UserRepository
	jwtKey []byte
}

// NewAuthService khởi tạo service
func NewAuthService(repo repositories.UserRepository, jwtKey []byte) AuthService {
	return &authService{repo: repo, jwtKey: jwtKey}
}

func validate(username, password string) error {
	if n := len(username); n < 3 || n > 32 {
		return fmt.Errorf("%w: username phải 3-32 ký tự", ErrInvalidInput)
	}
	if n := len(password); n < 6 || n > 72 {
		return fmt.Errorf("%w: password phải 6-72 byte", ErrInvalidInput)
	}
	return nil
}

func (s *authService) Register(username, password string) error {
	// 1. Giới hạn độ dài username và password
	username = strings.TrimSpace(username)
	if err := validate(username, password); err != nil {
		return err
	}

	// 2. Băm mật khẩu
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. Lưu vào DB
	err = s.repo.CreateUser(&models.User{Username: username, Password: string(hash)})
	if errors.Is(err, repositories.ErrDuplicate) {
		return ErrUserExists
	}
	return err
}

func (s *authService) Login(username, password string) (string, error) {
	// 1. Tìm user trong DB
	user, err := s.repo.GetUserByUsername(strings.TrimSpace(username))
	if errors.Is(err, repositories.ErrNotFound) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	// 2. So sánh mật khẩu
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	// 3. Tạo JWT Token
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtKey)
}
func (s *authService) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return s.jwtKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
