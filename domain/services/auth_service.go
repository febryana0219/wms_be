package services

import (
	"errors"
	"os"
	"strconv"
	"time"
	"wms-be/domain/repository"
	"wms-be/infrastructure/jwt"

	"golang.org/x/crypto/bcrypt"
)

// AuthService interface
type AuthService interface {
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(refreshToken string) error
}

type authService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

// Konstruktor
func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return &authService{userRepo, refreshTokenRepo}
}

// getEnvDuration membaca TTL dari .env
func getEnvDuration(key string, defaultVal int64) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return time.Duration(defaultVal) * time.Second
	}

	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return time.Duration(defaultVal) * time.Second
	}
	return time.Duration(val) * time.Second
}

// ===================== LOGIN =====================
func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials: user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials: wrong password")
	}

	accessTTL := getEnvDuration("JWT_ACCESS_EXPIRE", 3600)
	refreshTTL := getEnvDuration("JWT_REFRESH_EXPIRE", 604800)

	accessToken, err := jwt.GenerateToken(user.ID.String(), time.Now().Add(accessTTL))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateToken(user.ID.String(), time.Now().Add(refreshTTL))
	if err != nil {
		return "", "", err
	}

	// simpan refresh token di DB (pakai method Store(userID, token, expiresAt))
	if err := s.refreshTokenRepo.Store(user.ID, refreshToken, time.Now().Add(refreshTTL)); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ===================== REFRESH TOKEN =====================
func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	// validasi JWT
	_, err := jwt.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	// cek di DB
	rt, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return "", "", errors.New("refresh token not found or revoked/expired")
	}

	accessTTL := getEnvDuration("JWT_ACCESS_EXPIRE", 3600)
	refreshTTL := getEnvDuration("JWT_REFRESH_EXPIRE", 604800)

	accessToken, err := jwt.GenerateToken(rt.UserID.String(), time.Now().Add(accessTTL))
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := jwt.GenerateToken(rt.UserID.String(), time.Now().Add(refreshTTL))
	if err != nil {
		return "", "", err
	}

	// revoke token lama
	if err := s.refreshTokenRepo.Revoke(refreshToken); err != nil {
		return "", "", err
	}

	// simpan token baru
	if err := s.refreshTokenRepo.Store(rt.UserID, newRefreshToken, time.Now().Add(refreshTTL)); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// ===================== LOGOUT =====================
func (s *authService) Logout(refreshToken string) error {
	return s.refreshTokenRepo.Revoke(refreshToken)
}
