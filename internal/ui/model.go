package ui

import (
	"reddit-tui/internal/data"
	"reddit-tui/internal/icons"
	"reddit-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	SidebarItems  []string
	Posts         []models.Post
	SidebarCursor int
	PostsCursor   int
	ActivePane    string
	Width         int
	Height        int
	PostsScroll   int
	PreviewScroll int
}

func InitialModel() Model {
	posts, err := data.LoadSamplePosts()
	if err != nil {
		posts = []models.Post{}
	}

	return Model{
		SidebarItems: []string{
			icons.Home + " Home",
			icons.Popular + " Popular",
			icons.Explore + " Explore",
			icons.Settings + " Settings",
			icons.Login + " Login/Auth",
		},
		Posts:         posts,
		SidebarCursor: 0,
		PostsCursor:   0,
		ActivePane:    "sidebar",
		Width:         80,
		Height:        24,
		PostsScroll:   0,
		PreviewScroll: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			switch m.ActivePane {
			case "sidebar":
				m.ActivePane = "posts"
			case "posts":
				m.ActivePane = "preview"
			case "preview":
				m.ActivePane = "sidebar"
			}
		case "u":
			// Upvote the current post (only in preview pane)
			if m.ActivePane == "preview" && m.PostsCursor >= 0 && m.PostsCursor < len(m.Posts) {
				m.Posts[m.PostsCursor].ToggleUpvote()
			}
		case "d":
			// Downvote the current post (only in preview pane)
			if m.ActivePane == "preview" && m.PostsCursor >= 0 && m.PostsCursor < len(m.Posts) {
				m.Posts[m.PostsCursor].ToggleDownvote()
			}
		case "up", "k":
			if m.ActivePane == "sidebar" {
				if m.SidebarCursor > 0 {
					m.SidebarCursor--
				}
			} else if m.ActivePane == "posts" {
				if m.PostsCursor > 0 {
					m.PostsCursor--
					m.PreviewScroll = 0 // Reset preview scroll when changing posts
					// Keep cursor in view
					if m.PostsCursor < m.PostsScroll {
						m.PostsScroll = m.PostsCursor
					}
				}
			} else if m.ActivePane == "preview" {
				if m.PreviewScroll > 0 {
					m.PreviewScroll--
				}
			}
		case "down", "j":
			if m.ActivePane == "sidebar" {
				if m.SidebarCursor < len(m.SidebarItems)-1 {
					m.SidebarCursor++
				}
			} else if m.ActivePane == "posts" {
				if m.PostsCursor < len(m.Posts)-1 {
					m.PostsCursor++
					m.PreviewScroll = 0
					visiblePosts := (m.Height - 3 - 4) / 4
					if visiblePosts < 1 {
						visiblePosts = 1
					}
					if m.PostsCursor >= m.PostsScroll+visiblePosts {
						m.PostsScroll = m.PostsCursor - visiblePosts + 1
					}
				}
			} else if m.ActivePane == "preview" {
				m.PreviewScroll++
			}
		}
	}

	return m, nil
}
