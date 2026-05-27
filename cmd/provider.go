package cmd

import (
	"github.com/CyberGrit/go-spotify-me/internal/spotify"
	tea "github.com/charmbracelet/bubbletea"
)

type DataProvider interface {
	FetchTopArtists(timeRange string) tea.Cmd
	FetchTopSongs(timeRange string) tea.Cmd
	FetchArtistsPage(url string) tea.Cmd
	FetchSongsPage(url string) tea.Cmd
}

type DefaultDataProvider struct {
	client spotify.Client
}

func (d DefaultDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	return func() tea.Msg {
		url := "https://api.spotify.com/v1/me/top/artists?time_range=" + timeRange
		response, err := fetchArtistsPage(d.client, url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (d DefaultDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	return func() tea.Msg {
		url := "https://api.spotify.com/v1/me/top/tracks?time_range=" + timeRange
		response, err := fetchSongsPage(d.client, url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}

func (d DefaultDataProvider) FetchArtistsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchArtistsPage(d.client, url)
		if err != nil {
			return errMsg{err}
		}
		return switchToArtistsMsg{response}
	}
}

func (d DefaultDataProvider) FetchSongsPage(url string) tea.Cmd {
	return func() tea.Msg {
		response, err := fetchSongsPage(d.client, url)
		if err != nil {
			return errMsg{err}
		}
		return switchToSongsMsg{response}
	}
}
