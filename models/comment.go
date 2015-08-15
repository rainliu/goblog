package models

import (
	"time"

	"appengine/user"
)

type Comment struct {
	CommentID int
	ArticleID int
	Date      time.Time
	Name      string
	Email     string
	Website   string
	Content   string
}

type SideBar struct {
	User           *user.User
	LoginoutURL    string
	RSSConfig      BlogConfig
	RecentComments []Comment
	RecentPosts    []Article
	PopularPosts   []Article
	Archives       []string
	Links          []string
}

func CommentContent2Excerpt(text string) string {
	if len(text) < MyBlogConfig.ExcerptsCommentsCharNumber {
		return text
	} else {
		return text[0:MyBlogConfig.ExcerptsCommentsCharNumber] + " [...]"
	}
}
