package ui

import (
	"fmt"
	"strings"

	"reddit-tui/internal/models"

	"github.com/charmbracelet/lipgloss"
)

func renderPane(content string, width, height int, borderColor string, active bool) string {
	innerWidth := width - 2
	innerHeight := height - 2

	if innerWidth < 1 {
		innerWidth = 1
	}
	if innerHeight < 1 {
		innerHeight = 1
	}

	lines := strings.Split(content, "\n")
	result := make([]string, innerHeight)

	for i := 0; i < innerHeight; i++ {
		if i < len(lines) {
			line := lines[i]
			w := lipgloss.Width(line)
			if w > innerWidth {
				runes := []rune(line)
				if len(runes) > innerWidth {
					line = string(runes[:innerWidth])
				}
			}
			result[i] = line + strings.Repeat(" ", max(0, innerWidth-lipgloss.Width(line)))
		} else {
			result[i] = strings.Repeat(" ", innerWidth)
		}
	}

	innerContent := strings.Join(result, "\n")

	color := lipgloss.Color(borderColor)
	if active {
		color = lipgloss.Color("#ff5700")
	}

	style := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight)

	if borderColor != "" {
		style = style.BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(color)
	}

	return style.Render(innerContent)
}

func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return ""
	}

	controlPaneHeight := 3

	sidebarWidth := m.Width / 5
	if sidebarWidth < 15 {
		sidebarWidth = 15
	}
	remainingWidth := m.Width - sidebarWidth

	// Adjust layout based on whether settings is shown
	var postsWidth, previewWidth int
	if m.ShowSettings {
		// Full width for settings, no preview pane
		postsWidth = remainingWidth
		previewWidth = 0
	} else {
		// Normal layout with preview pane
		postsWidth = remainingWidth / 2
		previewWidth = remainingWidth - postsWidth
	}

	paneHeight := m.Height - controlPaneHeight

	postsPaneHeading := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4500")).Bold(true).MarginLeft(2)
	// previewPaneHeading := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4500")).Bold(true).MarginLeft(2)
	postTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb090"))
	postTitleSelectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5700"))
	subredditStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Center)
	previewTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5700")).MarginLeft(2)
	previewSubredditStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33")).MarginLeft(2)
	previewMetaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2)
	previewTextStyle := lipgloss.NewStyle().MarginLeft(2)

	sidebarItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffb090")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffb090")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(sidebarWidth - 4)
	sidebarItemActiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff5700")).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ff5700")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(sidebarWidth - 4)

	postItemStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ffb090")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(postsWidth - 6)
	postItemActiveStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ff5700")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(postsWidth - 6)

	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4500")).Bold(true)
	sidebarContent := logoStyle.Render(" ┬─┐┌─┐╔╦╗╦ ╦╦") + "\n"
	sidebarContent += logoStyle.Render(" ├┬┘├┤  ║ ║ ║║") + "\n"
	sidebarContent += logoStyle.Render(" ┴└─└─┘ ╩ ╚═╝╩") + "\n\n"
	for i, item := range m.SidebarItems {
		style := sidebarItemStyle
		if m.SidebarCursor == i {
			style = sidebarItemActiveStyle
		}
		sidebarContent += style.Render(item) + "\n"
	}

	var postsContent string

	if m.ShowSettings {
		// Settings pane
		postsContent = postsPaneHeading.Render("SETTINGS") + "\n\n"

		settingsLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true).MarginLeft(2)
		settingsInputStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ffb090")).
			PaddingLeft(1).
			PaddingRight(1).
			Width(postsWidth - 8)
		settingsInputActiveStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff5700")).
			PaddingLeft(1).
			PaddingRight(1).
			Width(postsWidth - 8)
		settingsHintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2)

		// API Key field
		apiKeyLabel := settingsLabelStyle.Render("API Key")
		apiKeyStyle := settingsInputStyle
		if m.SettingsCursor == 0 {
			apiKeyStyle = settingsInputActiveStyle
		}
		apiKeyValue := m.APIKey
		if m.EditingField == 1 {
			apiKeyValue += lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5700")).Render("█")
		}
		if apiKeyValue == "" && m.EditingField != 1 {
			apiKeyValue = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Enter your Reddit API key...")
		}

		postsContent += apiKeyLabel + "\n"
		postsContent += apiKeyStyle.Render(apiKeyValue) + "\n\n"

		// Client Secret field
		clientSecretLabel := settingsLabelStyle.Render("Client Secret")
		clientSecretStyle := settingsInputStyle
		if m.SettingsCursor == 1 {
			clientSecretStyle = settingsInputActiveStyle
		}
		clientSecretValue := m.ClientSecret
		// Mask the client secret
		if len(m.ClientSecret) > 0 && m.EditingField != 2 {
			clientSecretValue = strings.Repeat("•", len(m.ClientSecret))
		}
		if m.EditingField == 2 {
			clientSecretValue += lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5700")).Render("█")
		}
		if clientSecretValue == "" && m.EditingField != 2 {
			clientSecretValue = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Enter your client secret...")
		}

		postsContent += clientSecretLabel + "\n"
		postsContent += clientSecretStyle.Render(clientSecretValue) + "\n\n"

		postsContent += settingsHintStyle.Render("↑↓: navigate | Enter: edit | Esc: done") + "\n"

	} else if m.IsSearching {
		searchBarStyle := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ff5700")).
			PaddingLeft(1).
			PaddingRight(1).
			Width(postsWidth - 6)

		searchIconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4500")).Bold(true)
		searchBarContent := searchIconStyle.Render("Search: ") + m.SearchQuery
		if m.ActivePane == "posts" {
			searchBarContent += lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5700")).Render("█")
		}

		postsContent = postsPaneHeading.Render("EXPLORE") + "\n\n"
		postsContent += searchBarStyle.Render(searchBarContent) + "\n\n"

		if m.SearchQuery == "" {
			hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Center)
			postsContent += hintStyle.Render("Type to search posts...") + "\n"
		} else if len(m.SearchResults) == 0 {
			noResultsStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Align(lipgloss.Center)
			postsContent += noResultsStyle.Render("No results found") + "\n"
		} else {
			// Show search results
			// Calculate how many posts can fit (each post takes ~5 lines with border)
			// Account for extra space used by search bar (~4 lines)
			visiblePosts := (paneHeight - 8) / 5
			if visiblePosts < 1 {
				visiblePosts = 1
			}
			for i, post := range m.SearchResults {
				if i < m.PostsScroll {
					continue
				}
				if i >= m.PostsScroll+visiblePosts {
					break
				}
				titleStyle := postTitleStyle
				itemStyle := postItemStyle
				if m.PostsCursor == i {
					titleStyle = postTitleSelectedStyle
					itemStyle = postItemActiveStyle
				}

				postItemContent := titleStyle.Render(post.Title) + "\n"
				postItemContent += subredditStyle.Render(post.Subreddit) + " by u/" + post.Author + "\n"
				postItemContent += metaStyle.Render(fmt.Sprintf("%d upvotes | %d comments", post.GetDisplayUpvotes(), post.Comments))

				postsContent += itemStyle.Render(postItemContent) + "\n"
			}
		}
	} else {
		postsContent = postsPaneHeading.Render("POSTS") + "\n\n"
		// Calculate how many posts can fit (each post takes ~4 lines with border)
		visiblePosts := (paneHeight - 4) / 5 // Heading + spacing + post items
		if visiblePosts < 1 {
			visiblePosts = 1
		}
		for i, post := range m.Posts {
			if i < m.PostsScroll {
				continue
			}
			if i >= m.PostsScroll+visiblePosts {
				break
			}
			titleStyle := postTitleStyle
			itemStyle := postItemStyle
			if m.PostsCursor == i {
				titleStyle = postTitleSelectedStyle
				itemStyle = postItemActiveStyle
			}

			postItemContent := titleStyle.Render(post.Title) + "\n"
			postItemContent += subredditStyle.Render(post.Subreddit) + " by u/" + post.Author + "\n"
			postItemContent += metaStyle.Render(fmt.Sprintf("%d upvotes | %d comments", post.GetDisplayUpvotes(), post.Comments))

			postsContent += itemStyle.Render(postItemContent) + "\n"
		}
	}

	var previewLines []string

	// Determine which post to preview
	var selectedPost *models.Post
	if m.IsSearching && len(m.SearchResults) > 0 {
		if m.PostsCursor >= 0 && m.PostsCursor < len(m.SearchResults) {
			selectedPost = &m.SearchResults[m.PostsCursor]
		}
	} else if !m.IsSearching && len(m.Posts) > 0 {
		if m.PostsCursor >= 0 && m.PostsCursor < len(m.Posts) {
			selectedPost = &m.Posts[m.PostsCursor]
		}
	}

	if selectedPost != nil {

		// Vote indicators
		upvoteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2)
		downvoteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2)
		upvoteIcon := "▲"
		downvoteIcon := "▼"

		// Highlight active vote
		if selectedPost.UserVote == 1 { // VoteUp
			upvoteStyle = upvoteStyle.Foreground(lipgloss.Color("208")).Bold(true)
		} else if selectedPost.UserVote == 2 { // VoteDown
			downvoteStyle = downvoteStyle.Foreground(lipgloss.Color("33")).Bold(true)
		}

		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewTitleStyle.Render(selectedPost.Title))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewSubredditStyle.Render(selectedPost.Subreddit+" by u/"+selectedPost.Author))

		// Vote section
		voteCountStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true).MarginLeft(2)
		voteLine := upvoteStyle.Render(upvoteIcon) + " " + voteCountStyle.Render(fmt.Sprintf("%d", selectedPost.GetDisplayUpvotes())) + " " + downvoteStyle.Render(downvoteIcon)
		voteHintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(4)
		voteHint := voteHintStyle.Render("(u: upvote, d: downvote)")

		previewLines = append(previewLines, voteLine+" "+voteHint)
		previewLines = append(previewLines, previewMetaStyle.Render(fmt.Sprintf("%d comments", selectedPost.Comments)))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewTextStyle.Render(strings.Repeat("-", 20)))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewTextStyle.Render("Lorem ipsum dolor sit amet,"))
		previewLines = append(previewLines, previewTextStyle.Render("consectetur adipiscing elit."))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewTextStyle.Render("Sed do eiusmod tempor incididunt"))
		previewLines = append(previewLines, previewTextStyle.Render("ut labore et dolore magna aliqua."))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewTextStyle.Render("Ut enim ad minim veniam, quis"))
		previewLines = append(previewLines, previewTextStyle.Render("nostrud exercitation ullamco."))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewTextStyle.Render("Duis aute irure dolor in"))
		previewLines = append(previewLines, previewTextStyle.Render("reprehenderit in voluptate velit."))
	} else {
		previewLines = []string{"PREVIEW", "", "Select a post to view"}
	}

	scrollOffset := m.PreviewScroll
	if scrollOffset > len(previewLines)-1 {
		scrollOffset = len(previewLines) - 1
	}
	if scrollOffset < 0 {
		scrollOffset = 0
	}
	previewContent := strings.Join(previewLines[scrollOffset:], "\n")

	sidebar := renderPane(sidebarContent, sidebarWidth, paneHeight, "#ffb090", m.ActivePane == "sidebar")
	posts := renderPane(postsContent, postsWidth, paneHeight, "#ffb090", m.ActivePane == "posts")

	// Conditionally render main content based on settings view
	var mainContent string
	if m.ShowSettings {
		// No preview pane in settings
		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, posts)
	} else {
		// Include preview pane
		preview := renderPane(previewContent, previewWidth, paneHeight, "#ffb090", m.ActivePane == "preview")
		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, posts, preview)
	}

	controlTextStyle := metaStyle.Width(m.Width - 4)
	var controlText string
	if m.ShowSettings {
		controlText = controlTextStyle.Render("Enter: select section | Tab: switch panes | ↑↓/j/k: navigate | Esc: exit editing | q: quit")
	} else if m.IsSearching {
		controlText = controlTextStyle.Render("Enter: select section | Tab: switch panes | ↑↓/j/k: navigate | Esc: clear search | u: upvote | d: downvote | q: quit")
	} else {
		controlText = controlTextStyle.Render("Enter: select section | Tab: switch panes | ↑↓/j/k: navigate/scroll | u: upvote | d: downvote | q: quit")
	}
	controlPane := renderPane(controlText, m.Width, controlPaneHeight, "", false)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, controlPane)
}
