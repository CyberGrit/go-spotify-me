package cmd

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// mockDataProvider is a fake implementation of DataProvider for testing.
type mockDataProvider struct {
	lastMethodCalled string
	lastArg          string
}

func (m *mockDataProvider) FetchTopArtists(timeRange string) tea.Cmd {
	m.lastMethodCalled = "FetchTopArtists"
	m.lastArg = timeRange
	return func() tea.Msg {
		return switchToArtistsMsg{
			response: APIResponse{
				Artists: []Artist{{Name: "Mock Artist"}},
			},
		}
	}
}

func (m *mockDataProvider) FetchTopSongs(timeRange string) tea.Cmd {
	m.lastMethodCalled = "FetchTopSongs"
	m.lastArg = timeRange
	return func() tea.Msg {
		return switchToSongsMsg{
			response: APIResponse{
				Songs: []Song{{Name: "Mock Song"}},
			},
		}
	}
}

func (m *mockDataProvider) FetchArtistsPage(url string) tea.Cmd {
	m.lastMethodCalled = "FetchArtistsPage"
	m.lastArg = url
	return func() tea.Msg {
		return switchToArtistsMsg{}
	}
}

func (m *mockDataProvider) FetchSongsPage(url string) tea.Cmd {
	m.lastMethodCalled = "FetchSongsPage"
	m.lastArg = url
	return func() tea.Msg {
		return switchToSongsMsg{}
	}
}

func TestAppModelUpdate_KeyBindings(t *testing.T) {
	mockProvider := &mockDataProvider{}
	m := appModel{
		currentView: viewMenu,
		provider:    mockProvider,
		artistTable: table.New(),
		songTable:   table.New(),
	}

	tests := []struct {
		name           string
		keyMsg         string
		initialView    viewType
		expectedMethod string
		expectedArg    string
	}{
		{
			name:           "Press 'a' in menu view",
			keyMsg:         "a",
			initialView:    viewMenu,
			expectedMethod: "FetchTopArtists",
			expectedArg:    "medium_term",
		},
		{
			name:           "Press 's' in menu view",
			keyMsg:         "s",
			initialView:    viewMenu,
			expectedMethod: "FetchTopSongs",
			expectedArg:    "medium_term",
		},
		{
			name:           "Press '1' in artists view (short_term)",
			keyMsg:         "1",
			initialView:    viewArtists,
			expectedMethod: "FetchTopArtists",
			expectedArg:    "short_term",
		},
		{
			name:           "Press '3' in songs view (long_term)",
			keyMsg:         "3",
			initialView:    viewSongs,
			expectedMethod: "FetchTopSongs",
			expectedArg:    "long_term",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.currentView = tt.initialView
			mockProvider.lastMethodCalled = ""
			mockProvider.lastArg = ""

			// We handle '1', '3' vs 'a' where 'a' is length 1.
			// Let's create proper KeyMsg.
			var keyMsg tea.KeyMsg
			switch tt.keyMsg {
			case "1", "2", "3", "a", "s":
				keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.keyMsg)}
			}

			_, cmd := m.Update(keyMsg)

			if cmd == nil {
				t.Fatalf("expected tea.Cmd to be returned, got nil")
			}

			if mockProvider.lastMethodCalled != tt.expectedMethod {
				t.Errorf("expected method %s, got %s", tt.expectedMethod, mockProvider.lastMethodCalled)
			}

			if mockProvider.lastArg != tt.expectedArg {
				t.Errorf("expected arg %s, got %s", tt.expectedArg, mockProvider.lastArg)
			}
			
			// Test executing the command
			resMsg := cmd()
			
			switch tt.expectedMethod {
			case "FetchTopArtists", "FetchArtistsPage":
				if _, ok := resMsg.(switchToArtistsMsg); !ok {
					t.Errorf("expected switchToArtistsMsg, got %T", resMsg)
				}
			case "FetchTopSongs", "FetchSongsPage":
				if _, ok := resMsg.(switchToSongsMsg); !ok {
					t.Errorf("expected switchToSongsMsg, got %T", resMsg)
				}
			}
		})
	}
}

func TestAppModelUpdate_Pagination(t *testing.T) {
	mockProvider := &mockDataProvider{}
	m := appModel{
		currentView: viewArtists,
		provider:    mockProvider,
		artists: APIResponse{
			Next: "http://next-artists",
			Prev: "http://prev-artists",
		},
		songs: APIResponse{
			Next: "http://next-songs",
			Prev: "http://prev-songs",
		},
	}

	tests := []struct {
		name           string
		keyMsg         tea.KeyType
		initialView    viewType
		expectedMethod string
		expectedArg    string
	}{
		{
			name:           "Right arrow in artists view",
			keyMsg:         tea.KeyRight,
			initialView:    viewArtists,
			expectedMethod: "FetchArtistsPage",
			expectedArg:    "http://next-artists",
		},
		{
			name:           "Left arrow in artists view",
			keyMsg:         tea.KeyLeft,
			initialView:    viewArtists,
			expectedMethod: "FetchArtistsPage",
			expectedArg:    "http://prev-artists",
		},
		{
			name:           "Right arrow in songs view",
			keyMsg:         tea.KeyRight,
			initialView:    viewSongs,
			expectedMethod: "FetchSongsPage",
			expectedArg:    "http://next-songs",
		},
		{
			name:           "Left arrow in songs view",
			keyMsg:         tea.KeyLeft,
			initialView:    viewSongs,
			expectedMethod: "FetchSongsPage",
			expectedArg:    "http://prev-songs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.currentView = tt.initialView
			mockProvider.lastMethodCalled = ""
			mockProvider.lastArg = ""

			keyMsg := tea.KeyMsg{Type: tt.keyMsg}

			_, cmd := m.Update(keyMsg)

			if cmd == nil {
				t.Fatalf("expected tea.Cmd to be returned, got nil")
			}

			if mockProvider.lastMethodCalled != tt.expectedMethod {
				t.Errorf("expected method %s, got %s", tt.expectedMethod, mockProvider.lastMethodCalled)
			}

			if mockProvider.lastArg != tt.expectedArg {
				t.Errorf("expected arg %s, got %s", tt.expectedArg, mockProvider.lastArg)
			}
		})
	}
}
