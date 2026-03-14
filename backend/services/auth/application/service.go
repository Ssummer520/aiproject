package application

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"travel-api/services/auth/domain"
	authInfra "travel-api/services/auth/infrastructure"
)

var (
	ErrEmailExists  = errors.New("email already registered")
	ErrInvalidCreds = errors.New("invalid email or password")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidToken = errors.New("invalid or expired token")
)

type AuthService struct {
	userRepo   *authInfra.FileUserRepo
	tokenStore *authInfra.FileTokenStore
	resetStore *authInfra.FileTokenStore
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:   authInfra.NewFileUserRepo(),
		tokenStore: authInfra.NewFileTokenStore(authInfra.TokensFile, authInfra.TokenTTL),
		resetStore: authInfra.NewFileTokenStore(authInfra.ResetTokensFile, authInfra.ResetTokenTTL),
	}
}

const salt = "travel-api-auth-v1"

func hashPassword(password string) string {
	h := sha256.Sum256([]byte(salt + password))
	return hex.EncodeToString(h[:])
}

func checkPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generateID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("user_%x", b), nil
}

func (s *AuthService) Register(email, password string) (*domain.User, error) {
	if _, ok := s.userRepo.FindByEmail(email); ok {
		return nil, ErrEmailExists
	}
	id, err := generateID()
	if err != nil {
		return nil, err
	}
	u := domain.User{
		ID:           id,
		Email:        email,
		PasswordHash: hashPassword(password),
		CreatedAt:    time.Now(),
	}
	if err := s.userRepo.Create(u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *AuthService) Login(email, password string) (token string, user *domain.User, err error) {
	u, ok := s.userRepo.FindByEmail(email)
	if !ok || !checkPassword(password, u.PasswordHash) {
		return "", nil, ErrInvalidCreds
	}
	token, err = generateToken()
	if err != nil {
		return "", nil, err
	}
	s.tokenStore.Set(token, u.ID)
	return token, u, nil
}

func (s *AuthService) ForgotPassword(email string) (resetToken string, err error) {
	u, ok := s.userRepo.FindByEmail(email)
	if !ok {
		return "", ErrUserNotFound
	}
	resetToken, err = generateToken()
	if err != nil {
		return "", err
	}
	s.resetStore.Set(resetToken, u.ID)
	return resetToken, nil
}

func (s *AuthService) ResetPassword(resetToken, newPassword string) error {
	userID, ok := s.resetStore.Get(resetToken)
	if !ok {
		return ErrInvalidToken
	}
	u, ok := s.userRepo.FindByID(userID)
	if !ok {
		return ErrUserNotFound
	}
	u.PasswordHash = hashPassword(newPassword)
	if err := s.userRepo.Update(u); err != nil {
		return err
	}
	s.resetStore.Delete(resetToken)
	return nil
}

func (s *AuthService) ValidateToken(token string) (userID string, ok bool) {
	return s.tokenStore.Get(token)
}

func (s *AuthService) Logout(token string) {
	s.tokenStore.Delete(token)
}

func (s *AuthService) GetUserByID(userID string) (*domain.User, bool) {
	return s.userRepo.FindByID(userID)
}
