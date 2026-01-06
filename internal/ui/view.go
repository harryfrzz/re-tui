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
	postsWidth := remainingWidth / 2
	previewWidth := remainingWidth - postsWidth

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
	sidebarContent := logoStyle.Render(" ┌─┐┌─┐┌┬┐┌┬┐┬┌┬┐") + "\n"
	sidebarContent += logoStyle.Render(" ├┬┘├┤  ││ ││││││") + "\n"
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

	var previewLines []string
	if m.PostsCursor >= 0 && m.PostsCursor < len(m.Posts) {
		selectedPost := m.Posts[m.PostsCursor]

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
	preview := renderPane(previewContent, previewWidth, paneHeight, "#ffb090", m.ActivePane == "preview")

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, posts, preview)

	controlTextStyle := metaStyle.Width(m.Width - 4)
	controlText := controlTextStyle.Render("Tab: switch panes | ↑↓/j/k: navigate/scroll | u: upvote | d: downvote | q: quit")
	controlPane := renderPane(controlText, m.Width, controlPaneHeight, "", false)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, controlPane)
}
