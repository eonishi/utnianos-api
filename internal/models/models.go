package models

import "time"

type ForumLastPost struct {
	URL       string    `json:"url,omitempty"`
	Date      *time.Time `json:"date,omitempty"`
	Author    string    `json:"author,omitempty"`
	AuthorURL string    `json:"author_url,omitempty"`
}
type Forum struct {
	Name         string        `json:"name"`
	URL          string        `json:"url"`
	Slug         string        `json:"slug,omitempty"`
	TopicCount   int           `json:"topic_count,omitempty"`
	MessageCount int           `json:"message_count,omitempty"`
	LastPost     *ForumLastPost `json:"last_post,omitempty"` // Es posible que esto sea last topic por como esta armado la API
}

type Topic struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	URL          string        `json:"url"`
	Author       string        `json:"author"`
	AuthorURL    string        `json:"author_url,omitempty"`
	Replies      int           `json:"replies,omitempty"`
	Views        int           `json:"views,omitempty"`
	ThankedCount int           `json:"thanked_count,omitempty"`
	Pages        int           `json:"pages,omitempty"`
	Materias     []string      `json:"materias,omitempty"`
	Aportes      []string      `json:"tipo_aportes,omitempty"`
	LastPost     *ForumLastPost `json:"last_post,omitempty"`
}

type TopicDetail struct {
	Topic
	Posts []Post `json:"posts"`
}

type Post struct {
	ID           int       `json:"id"`
	Author       string    `json:"author"`
	AuthorURL    string    `json:"author_url"`
	Content      string    `json:"content"`
	Date         time.Time `json:"date"`
	Links        []string  `json:"links,omitempty"`
	Attachments []string  `json:"attachments,omitempty"`
	ReplyPostID  int       `json:"reply_post_id,omitempty"`
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	JoinDate  string `json:"join_date"`
	Posts     int    `json:"posts"`
	Location  string `json:"location"`
	Website   string `json:"website,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type SearchResult struct {
	Topics     []Topic `json:"topics"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	TotalPages int     `json:"total_pages"`
}

type ForumListResponse struct {
	Forums []Forum `json:"forums"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
