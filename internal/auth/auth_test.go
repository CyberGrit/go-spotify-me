package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"regexp"
	"testing"
)

func TestGenerateCodeVerifier(t *testing.T) {
	verifier := GenerateCodeVerifier()
	
	if len(verifier) != 64 {
		t.Errorf("Expected length 64, got %d", len(verifier))
	}

	match, _ := regexp.MatchString("^[a-z]{64}$", verifier)
	if !match {
		t.Errorf("Expected verifier to contain only a-z characters, got %s", verifier)
	}
}

func TestGenerateCodeChallenge(t *testing.T) {
	verifier := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmn"
	
	expectedHash := sha256.Sum256([]byte(verifier))
	expectedChallenge := base64.RawURLEncoding.EncodeToString(expectedHash[:])
	
	challenge := GenerateCodeChallenge(verifier)
	
	if challenge != expectedChallenge {
		t.Errorf("Expected challenge %s, got %s", expectedChallenge, challenge)
	}
}
