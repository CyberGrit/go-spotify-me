package auth

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zalando/go-keyring"
	"go.uber.org/zap"
)

// TokenStore defines the interface for token storage.
type TokenStore interface {
	GetRefreshToken() (string, error)
	SetRefreshToken(token string) error
	
	GetClientID() (string, error)
	SetClientID(id string) error
	
	GetAccessToken() (string, time.Time, error)
	SetAccessToken(token string, expiresAt time.Time, refreshToken string) error
	
	Clear() error
}

// OSStore implements TokenStore using OS keyring and filesystem.
type OSStore struct {
	serviceName string
	fileName    string
	logger      *zap.Logger
}

// NewOSStore creates a new OSStore.
func NewOSStore(logger *zap.Logger) *OSStore {
	return &OSStore{
		serviceName: "go-spotify-me-cli",
		fileName:    ".go-spotify-me-cli",
		logger:      logger,
	}
}

func (s *OSStore) getFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, s.fileName)

	// Validate that the filePath is within the user's home directory
	if !strings.HasPrefix(filePath, homeDir) {
		return "", fmt.Errorf("invalid file path: %s", filePath)
	}

	return filePath, nil
}

func (s *OSStore) GetRefreshToken() (string, error) {
	// Check for refresh token in the keyring
	refreshToken, err := keyring.Get(s.serviceName, "refresh_token")
	if err == nil && refreshToken != "" {
		return refreshToken, nil
	}

	if s.logger != nil {
		s.logger.Debug("Refresh token not found in keyring", zap.Error(err))
	}

	// Check for refresh token in the hidden file
	filePath, err := s.getFilePath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "refresh_token=") {
			return strings.TrimPrefix(line, "refresh_token="), nil
		}
	}

	return "", fmt.Errorf("refresh token not found")
}

func (s *OSStore) SetRefreshToken(token string) error {
	// Setting the refresh token is usually handled by SetAccessToken since they are often updated together.
	return keyring.Set(s.serviceName, "refresh_token", token)
}

func (s *OSStore) GetClientID() (string, error) {
	clientID, err := keyring.Get(s.serviceName, "client_id")
	if err == nil && clientID != "" {
		return clientID, nil
	}
	return "", fmt.Errorf("client ID not found in keyring")
}

func (s *OSStore) SetClientID(id string) error {
	return keyring.Set(s.serviceName, "client_id", id)
}

func (s *OSStore) GetAccessToken() (string, time.Time, error) {
	filePath, err := s.getFilePath()
	if err != nil {
		return "", time.Time{}, err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to read token file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var accessToken, expiresAtStr string
	for _, line := range lines {
		if strings.HasPrefix(line, "access_token=") {
			accessToken = strings.TrimPrefix(line, "access_token=")
		} else if strings.HasPrefix(line, "expires_at=") {
			expiresAtStr = strings.TrimPrefix(line, "expires_at=")
		}
	}

	if accessToken == "" || expiresAtStr == "" {
		return "", time.Time{}, fmt.Errorf("access token or expiration not found")
	}

	expirationTime, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse expiration time: %w", err)
	}

	if time.Now().After(expirationTime) {
		return "", time.Time{}, fmt.Errorf("access token is expired")
	}

	return accessToken, expirationTime, nil
}

func (s *OSStore) SetAccessToken(token string, expiresAt time.Time, refreshToken string) error {
	filePath, err := s.getFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	var data string

	// Attempt to store the refresh token in the keyring
	if refreshToken != "" {
		err = keyring.Set(s.serviceName, "refresh_token", refreshToken)
		if err != nil {
			if s.logger != nil {
				s.logger.Error("Failed to store refresh token in keyring", zap.Error(err))
				s.logger.Info("Falling back to saving the refresh token in the hidden file.")
			}
			data = fmt.Sprintf("access_token=%s\nrefresh_token=%s\nexpires_at=%s\n", token, refreshToken, expiresAt.Format(time.RFC3339))
		} else {
			data = fmt.Sprintf("access_token=%s\nexpires_at=%s\n", token, expiresAt.Format(time.RFC3339))
		}
	} else {
		data = fmt.Sprintf("access_token=%s\nexpires_at=%s\n", token, expiresAt.Format(time.RFC3339))
	}

	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	if s.logger != nil {
		s.logger.Debug("Access token saved to file", zap.String("filePath", filePath))
	}

	return nil
}

func (s *OSStore) Clear() error {
	// Remove client_id from the keyring
	if err := keyring.Delete(s.serviceName, "client_id"); err != nil {
		if s.logger != nil {
			s.logger.Debug("Failed to delete client_id from keyring", zap.Error(err))
		}
	}

	// Remove refresh_token from the keyring
	if err := keyring.Delete(s.serviceName, "refresh_token"); err != nil {
		if s.logger != nil {
			s.logger.Debug("Failed to delete refresh_token from keyring", zap.Error(err))
		}
	}

	// Clear the hidden file
	filePath, err := s.getFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete configuration file: %w", err)
	}

	return nil
}

// InMemoryStore implements TokenStore for testing purposes.
type InMemoryStore struct {
	refreshToken string
	clientID     string
	accessToken  string
	expiresAt    time.Time
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{}
}

func (s *InMemoryStore) GetRefreshToken() (string, error) {
	if s.refreshToken == "" {
		return "", fmt.Errorf("refresh token not found")
	}
	return s.refreshToken, nil
}

func (s *InMemoryStore) SetRefreshToken(token string) error {
	s.refreshToken = token
	return nil
}

func (s *InMemoryStore) GetClientID() (string, error) {
	if s.clientID == "" {
		return "", fmt.Errorf("client ID not found")
	}
	return s.clientID, nil
}

func (s *InMemoryStore) SetClientID(id string) error {
	s.clientID = id
	return nil
}

func (s *InMemoryStore) GetAccessToken() (string, time.Time, error) {
	if s.accessToken == "" {
		return "", time.Time{}, fmt.Errorf("access token not found")
	}
	if time.Now().After(s.expiresAt) {
		return "", time.Time{}, fmt.Errorf("access token is expired")
	}
	return s.accessToken, s.expiresAt, nil
}

func (s *InMemoryStore) SetAccessToken(token string, expiresAt time.Time, refreshToken string) error {
	s.accessToken = token
	s.expiresAt = expiresAt
	if refreshToken != "" {
		s.refreshToken = refreshToken
	}
	return nil
}

func (s *InMemoryStore) Clear() error {
	s.refreshToken = ""
	s.clientID = ""
	s.accessToken = ""
	s.expiresAt = time.Time{}
	return nil
}
