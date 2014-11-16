package store

import "time"

type User struct {
	ID         int
	Name       string
	Email      string
	Password   string // bcrypted
	CreatedAt  time.Time
	Management bool
}

type Blog struct {
	ID    int
	Name  string
	About string `type:"text"`
}

type BlogAuthor struct {
	ID     int
	Admin  bool
	BlogID int
	UserID int
}

type Post struct {
	ID        int
	Title     string `type:"text"`
	Permalink string
	Text      string `type:"text"` // markdown version of post
	CreatedAt time.Time
	PostedAt  time.Time

	UserID int // author
	BlogID int // network/blog

}

type Comment struct {
	ID   int
	Text string `type:"text"`

	UserID int
	PostID int
}
