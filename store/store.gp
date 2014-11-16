package store

type User table {
  ID int
  Name string
  Email string
  Password string // bcrypted
  CreatedAt time.Time
  Management bool

  relation {
    []BlogAuthor
    []Post
    []Comment
  }
}

type Blog table {
  ID int
  Name string
  About string `type:"text"`

  relation {
    []BlogAuthor
    []Post
  }
}

type BlogAuthor table {
  ID int
  Admin bool
  BlogID int
  UserID int

  relation {
    User
    Blog
  }
}

type Post table {
  ID int
  Title string `type:"text"`
  Permalink string
  Text string `type:"text"` // markdown version of post
  CreatedAt time.Time
  PostedAt time.Time

  UserID int // author
  BlogID int // network/blog

  relation {
    Blog
    User
    []Comment
  }
}

type Comment table {
  ID int
  Text string `type:"text"`

  UserID int
  PostID int

  relation {
    Post
    User
  }
}
