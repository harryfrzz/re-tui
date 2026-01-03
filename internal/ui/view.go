package ui

import (
	"fmt"
	"strings"

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

			if lipgloss.Width(line) > innerWidth {
				runes := []rune(line)
				if len(runes) > innerWidth {
					line = string(runes[:innerWidth])
				}
			}

			result[i] = line + strings.Repeat(" ", innerWidth-lipgloss.Width(line))
		} else {
			result[i] = strings.Repeat(" ", innerWidth)
		}
	}

	innerContent := strings.Join(result, "\n")

	color := lipgloss.Color(borderColor)
	if active {
		color = lipgloss.Color("205")
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Width(innerWidth).
		Height(innerHeight)

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
	postsWidth := remainingWidth / 2
	previewWidth := remainingWidth - postsWidth

	paneHeight := m.Height - controlPaneHeight

	postsPaneHeading := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginLeft(2)
	previewPaneHeading := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginLeft(2)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	postTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	subredditStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	previewTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginLeft(2)
	previewSubredditStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33")).MarginLeft(2)
	previewMetaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginLeft(2)
	previewTextStyle := lipgloss.NewStyle().MarginLeft(2)

	sidebarItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(sidebarWidth - 4)
	sidebarItemActiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(sidebarWidth - 4)

	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
	sidebarContent := logoStyle.Render(" ┌─┐┌─┐┌┬┐┌┬┐┬┌┬┐") + "\n"
	sidebarContent += logoStyle.Render(" ├┬┘├┤  ││ ││││|│ ") + "\n"
	sidebarContent += logoStyle.Render(" ┴└─└─┘─┴┘─┴┘┴ ┴ ") + "\n"
	sidebarContent += logoStyle.Render(" T U I       ") + "\n\n"
	for i, item := range m.SidebarItems {
		style := sidebarItemStyle
		if m.SidebarCursor == i {
			style = sidebarItemActiveStyle
		}
		sidebarContent += style.Render(item) + "\n"
	}

	postsContent := postsPaneHeading.Render("POSTS") + "\n\n"
	for i, post := range m.Posts {
		if i < m.PostsScroll {
			continue
		}
		cursor := "  "
		if m.PostsCursor == i {
			cursor = cursorStyle.Render("> ")
		}
		postsContent += cursor + postTitleStyle.Render(post.Title) + "\n"
		postsContent += "   " + subredditStyle.Render(post.Subreddit) + " by u/" + post.Author + "\n"
		postsContent += "   " + metaStyle.Render(fmt.Sprintf("%d upvotes | %d comments", post.Upvotes, post.Comments)) + "\n\n"
	}

	var previewLines []string
	if m.PostsCursor >= 0 && m.PostsCursor < len(m.Posts) {
		selectedPost := m.Posts[m.PostsCursor]
		previewLines = append(previewLines, previewPaneHeading.Render("PREVIEW"))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewTitleStyle.Render(selectedPost.Title))
		previewLines = append(previewLines, "")
		previewLines = append(previewLines, previewSubredditStyle.Render(selectedPost.Subreddit+" by u/"+selectedPost.Author))
		previewLines = append(previewLines, previewMetaStyle.Render(fmt.Sprintf("%d upvotes | %d comments", selectedPost.Upvotes, selectedPost.Comments)))
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

	sidebar := renderPane(sidebarContent, sidebarWidth, paneHeight, "63", m.ActivePane == "sidebar")
	posts := renderPane(postsContent, postsWidth, paneHeight, "63", m.ActivePane == "posts")
	preview := renderPane(previewContent, previewWidth, paneHeight, "63", m.ActivePane == "preview")

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, posts, preview)

	controlText := metaStyle.Render("Tab: switch panes | ↑↓/j/k: navigate/scroll | q: quit")
	controlPane := renderPane(controlText, m.Width, controlPaneHeight, "63", false)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, controlPane)
}
