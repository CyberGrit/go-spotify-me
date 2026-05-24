package cmd

import (
	"github.com/CyberGrit/go-spotify-me/internal/auth"
)

// ClearConfig removes the client_id and refresh_token from the store
// and clears the .go-spotify-me-cli file in the user's home directory.
func ClearConfig() error {
	store := auth.NewOSStore(nil)
	return store.Clear()
}
