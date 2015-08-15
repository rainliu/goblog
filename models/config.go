package models

type BlogConfig struct {
	Name                       string
	Description                string
	Url                        string
	RssUrl                     string
	ExcerptsCharNumber         int
	ExcerptsNumber             int
	ExcerptsCommentsCharNumber int
	RecentCommentsNumber       int
	PopularPostsNumber         int
	RecentPostsNumber          int
}

var MyBlogConfig = BlogConfig{Name: "H265.net",
	Description:                "Witness the development of H.265",
	Url:                        "http://localhost:8080",
	RssUrl:                     "http://localhost:8080/rss",
	ExcerptsCharNumber:         1000,
	ExcerptsNumber:             4,
	ExcerptsCommentsCharNumber: 70,
	RecentCommentsNumber:       10,
	PopularPostsNumber:         10,
	RecentPostsNumber:          10}
