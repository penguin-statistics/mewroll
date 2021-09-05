package httpclient

import "time"

type MewUser struct {
	Background  string      `json:"background"`
	Username    string      `json:"username"`
	Name        string      `json:"name"`
	Avatar      string      `json:"avatar"`
	Type        int         `json:"type"`
	Profile     string      `json:"profile"`
	Id          string      `json:"id"`
}

type MewThoughtComment struct {
	ThoughtId       string       `json:"thought_id"`
	AuthorId        string       `json:"author_id"`
	Content         string       `json:"content"`
	NodeId          string       `json:"node_id"`
	CreatedAt       time.Time    `json:"created_at"`
	Index           int          `json:"index"`
	Id              string        `json:"id"`
	LikeCount       int           `json:"like_count"`
	CommentCount    int           `json:"comment_count"`
	Deleted         bool          `json:"deleted,omitempty"`
}

type MewCommentResponse struct {
	Objects struct {
		Users map[string]MewUser `json:"users"`
	} `json:"objects"`
	Entries []MewThoughtComment `json:"entries"`
}
