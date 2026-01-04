package models

type VoteState int

const (
	VoteNone VoteState = iota
	VoteUp
	VoteDown
)

type Post struct {
	Title      string    `json:"title"`
	Subreddit  string    `json:"subreddit"`
	Author     string    `json:"author"`
	Upvotes    int       `json:"upvotes"`
	Comments   int       `json:"comments"`
	UserVote   VoteState 
	VoteOffset int      
}

//returns the upvote count adjusted for user vote
func (p *Post) GetDisplayUpvotes() int {
	return p.Upvotes + p.VoteOffset
}

// toggles upvote state
func (p *Post) ToggleUpvote() {
	if p.UserVote == VoteUp {
		// Remove upvote
		p.UserVote = VoteNone
		p.VoteOffset = 0
	} else if p.UserVote == VoteDown {
		// Switch from downvote to upvote
		p.UserVote = VoteUp
		p.VoteOffset = 2 // +1 to cancel downvote, +1 for upvote
	} else {
		// Add upvote
		p.UserVote = VoteUp
		p.VoteOffset = 1
	}
}

// ToggleDownvote toggles downvote state
func (p *Post) ToggleDownvote() {
	if p.UserVote == VoteDown {
		// Remove downvote
		p.UserVote = VoteNone
		p.VoteOffset = 0
	} else if p.UserVote == VoteUp {
		// Switch from upvote to downvote
		p.UserVote = VoteDown
		p.VoteOffset = -2 // -1 to cancel upvote, -1 for downvote
	} else {
		// Add downvote
		p.UserVote = VoteDown
		p.VoteOffset = -1
	}
}
