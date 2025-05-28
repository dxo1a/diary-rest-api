package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"auth-service/internal/domain"
	"auth-service/internal/repository"
)

var (
	ErrUserExists    = errors.New("user already exists")
	ErrInvalidCreds  = errors.New("invalid credentials")
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenRevoked  = errors.New("refresh token revoked")
	ErrTokenExpired  = errors.New("refresh token expired")
)

type AuthService struct {
	users      *repository.UserRepo
	tokens     *repository.TokenRepo
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthService(u *repository.UserRepo, t *repository.TokenRepo, secret string, at, rt time.Duration) *AuthService {
	return &AuthService{
		users:      u,
		tokens:     t,
		jwtSecret:  []byte(secret),
		accessTTL:  at,
		refreshTTL: rt,
	}
}

func (s *AuthService) Register(username, password string) error {
	exists, err := s.users.FindByUsername(username)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrUserExists
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return s.users.Create(&domain.User{Username: username, PasswordHash: string(hash)})
}

func (s *AuthService) Login(username, password string) (access, refresh string, err error) {
	user, err := s.users.FindByUsername(username)
	if err != nil || user == nil {
		return "", "", ErrInvalidCreds
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", "", ErrInvalidCreds
	}
	// Access JWT
	atClaims := jwt.RegisteredClaims{
		Subject:   string(rune(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
	}
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	access, err = aToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}
	// Refresh token
	rt := domain.RefreshToken{
		UserID:    user.ID,
		Token:     generateRandomString(32),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.tokens.Create(&rt); err != nil {
		return "", "", err
	}
	return access, rt.Token, nil
}

func (s *AuthService) Refresh(oldToken string) (newAccess, newRefresh string, err error) {
	// найти запись по старому RT
	rt, err := s.tokens.Find(oldToken)
	if err != nil || rt == nil {
		return "", "", ErrTokenNotFound
	}
	if rt.Revoked {
		return "", "", ErrTokenRevoked
	}
	if time.Now().After(rt.ExpiresAt) {
		return "", "", ErrTokenExpired
	}

	// отозвать старый RT
	if err := s.tokens.Revoke(rt.ID); err != nil {
		return "", "", err
	}

	// сгенерировать новый Access‑JWT
	sub := strconv.FormatUint(uint64(rt.UserID), 10)
	atClaims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	newAccess, err = at.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// сгенерировать и сохранить новый Refresh‑token
	randToken := generateRandomString(32)
	newRT := &domain.RefreshToken{
		UserID:    rt.UserID,
		Token:     randToken,
		ExpiresAt: time.Now().Add(s.refreshTTL),
	}
	if err := s.tokens.Create(newRT); err != nil {
		return "", "", err
	}
	newRefresh = newRT.Token

	return newAccess, newRefresh, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	rt, err := s.tokens.Find(refreshToken)
	if err != nil || rt == nil {
		return ErrTokenNotFound
	}
	return s.tokens.Revoke(rt.ID)
}

// вспомогательная генерация случайной строки
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
