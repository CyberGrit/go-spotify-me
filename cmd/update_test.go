package cmd

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type MockDataProvider struct {
	CalledFetchTopArtists  bool
	CalledFetchTopSongs    bool
	CalledFetchArtistsPage bool
	CalledFetchSongsPage   bool
}

func (m *MockDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	m.CalledFetchTopArtists = true
	return func() tea.Msg {
		return switchToArtistsMsg{}
	}
}

func (m *MockDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	m.CalledFetchTopSongs = true
	return func() tea.Msg {
		return switchToSongsMsg{}
	}
}

func (m *MockDataProvider) FetchArtistsPage(url string) tea.Cmd {
	m.CalledFetchArtistsPage = true
	return func() tea.Msg {
		return switchToArtistsMsg{}
	}
}

func (m *MockDataProvider) FetchSongsPage(url string) tea.Cmd {
	m.CalledFetchSongsPage = true
	return func() tea.Msg {
		return switchToSongsMsg{}
	}
}

func TestUpdate_FetchTopArtists(t *testing.T) {
	mockProvider := &MockDataProvider{}
	m := appModel{
		currentView: viewMenu,
		provider:    mockProvider,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatalf("Expected non-nil cmd")
	}

	resultMsg := cmd()

	if !mockProvider.CalledFetchTopArtists {
		t.Errorf("Expected FetchTopArtists to be called")
	}

	if _, ok := resultMsg.(switchToArtistsMsg); !ok {
		t.Errorf("Expected switchToArtistsMsg, got %T", resultMsg)
	}
}

func TestUpdate_FetchTopSongs(t *testing.T) {
	mockProvider := &MockDataProvider{}
	m := appModel{
		currentView: viewMenu,
		provider:    mockProvider,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatalf("Expected non-nil cmd")
	}

	resultMsg := cmd()

	if !mockProvider.CalledFetchTopSongs {
		t.Errorf("Expected FetchTopSongs to be called")
	}

	if _, ok := resultMsg.(switchToSongsMsg); !ok {
		t.Errorf("Expected switchToSongsMsg, got %T", resultMsg)
	}
}

func TestUpdate_PaginationArtistsNext(t *testing.T) {
	mockProvider := &MockDataProvider{}
	m := appModel{
		currentView: viewArtists,
		artists: APIResponse{
			Next: "next-page-url",
		},
		provider:    mockProvider,
	}

	msg := tea.KeyMsg{Type: tea.KeyRight}
	_, cmd := m.Update(msg)

	if cmd == nil {
		t.Fatalf("Expected non-nil cmd")
	}

	resultMsg := cmd()

	if !mockProvider.CalledFetchArtistsPage {
		t.Errorf("Expected FetchArtistsPage to be called")
	}

	if _, ok := resultMsg.(switchToArtistsMsg); !ok {
		t.Errorf("Expected switchToArtistsMsg, got %T", resultMsg)
	}
}
