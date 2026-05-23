package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/CyberGrit/go-spotify-me/internal/auth"
)

type SpotifyClient interface {
	GetTopArtists(timeRange string) ([]Artist, string, string, error)
	GetTopSongs(timeRange string) ([]Song, string, string, error)
	GetMe() (Me, error)
	GetArtistsPage(url string) ([]Artist, string, string, error)
	GetSongsPage(url string) ([]Song, string, string, error)
}

type DefaultSpotifyClient struct{}

// makeAPIRequest makes a GET request to the Spotify API and returns the response or an error
func (c *DefaultSpotifyClient) makeAPIRequest(token string, url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	return response, nil
}

func (c *DefaultSpotifyClient) GetTopArtists(timeRange string) ([]Artist, string, string, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?time_range=%s&limit=50", timeRange)
	return c.GetArtistsPage(url)
}

func (c *DefaultSpotifyClient) GetArtistsPage(url string) ([]Artist, string, string, error) {
	token, _ := auth.GetValidAccessToken()
	response, err := c.makeAPIRequest(token, url)
	if err != nil {
		return nil, "", "", err
	}

	artists := parseArtists(response)
	next, _ := response["next"].(string)
	prev, _ := response["previous"].(string)

	return artists, next, prev, nil
}

func parseArtists(response map[string]interface{}) []Artist {
	items, ok := response["items"].([]interface{})
	if !ok {
		return nil
	}

	var artists []Artist
	for _, item := range items {
		artist, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name := artist["name"].(string)

		genresInterface, ok := artist["genres"].([]interface{})
		if !ok {
			genresInterface = []interface{}{}
		}
		var genres []string
		for _, genre := range genresInterface {
			genres = append(genres, genre.(string))
		}

		popularity := int(artist["popularity"].(float64))

		artists = append(artists, Artist{
			Name:       name,
			Genres:     strings.Join(genres, ", "),
			Popularity: popularity,
		})
	}

	return artists
}

func (c *DefaultSpotifyClient) GetTopSongs(timeRange string) ([]Song, string, string, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?time_range=%s&limit=50", timeRange)
	return c.GetSongsPage(url)
}

func (c *DefaultSpotifyClient) GetSongsPage(url string) ([]Song, string, string, error) {
	token, _ := auth.GetValidAccessToken()
	response, err := c.makeAPIRequest(token, url)
	if err != nil {
		return nil, "", "", err
	}

	songs := parseSongs(response)
	next, _ := response["next"].(string)
	prev, _ := response["previous"].(string)

	return songs, next, prev, nil
}

func parseSongs(response map[string]interface{}) []Song {
	items, ok := response["items"].([]interface{})
	if !ok {
		return nil
	}

	var songs []Song
	for _, item := range items {
		track, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		name := track["name"].(string)
		popularity := int(track["popularity"].(float64))

		albumName := ""
		if album, ok := track["album"].(map[string]interface{}); ok {
			albumName, _ = album["name"].(string)
		}

		artistName := ""
		if artistList, ok := track["artists"].([]interface{}); ok && len(artistList) > 0 {
			if firstArtist, ok := artistList[0].(map[string]interface{}); ok {
				artistName, _ = firstArtist["name"].(string)
			}
		}

		songs = append(songs, Song{
			Name:       name,
			Artist:     artistName,
			Album:      albumName,
			Popularity: popularity,
		})
	}

	return songs
}

func (c *DefaultSpotifyClient) GetMe() (Me, error) {
	token, _ := auth.GetValidAccessToken()
	response, err := c.makeAPIRequest(token, "https://api.spotify.com/v1/me")
	if err != nil {
		return Me{}, err
	}

	country, _ := response["country"].(string)
	displayName, _ := response["display_name"].(string)
	email, _ := response["email"].(string)
	product, _ := response["product"].(string)
	externalURLs, _ := response["external_urls"].(map[string]interface{})
	profileURL, _ := externalURLs["spotify"].(string)

	return Me{
		Country:     country,
		DisplayName: displayName,
		Email:       email,
		Product:     product,
		ProfileURL:  profileURL,
	}, nil
}
