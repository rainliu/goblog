package models

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

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

func GetSideBar(r *http.Request) (*SideBar, error) {
	var url string

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, _ = user.LoginURL(c, r.URL.String())
	} else {
		url, _ = user.LogoutURL(c, r.URL.String())
	}

	a := datastore.NewQuery("Comment").Order("-Date").Limit(MyBlogConfig.ExcerptsNumber)
	RecentComments := make([]Comment, 0, MyBlogConfig.ExcerptsNumber)
	if _, err := a.GetAll(c, &RecentComments); err != nil {
		return nil, err
	}

	b := datastore.NewQuery("Article").Order("-Date").Limit(MyBlogConfig.ExcerptsNumber)
	RecentPosts := make([]Article, 0, MyBlogConfig.ExcerptsNumber)
	if _, err := b.GetAll(c, &RecentPosts); err != nil {
		return nil, err
	}

	p := datastore.NewQuery("Article").Order("-Views").Limit(MyBlogConfig.ExcerptsNumber)
	PopularPosts := make([]Article, 0, MyBlogConfig.ExcerptsNumber)
	if _, err := p.GetAll(c, &PopularPosts); err != nil {
		return nil, err
	}

	mySideBar := &SideBar{
		User:           u,
		LoginoutURL:    url,
		RSSConfig:      MyBlogConfig,
		RecentComments: RecentComments,
		RecentPosts:    RecentPosts,
		PopularPosts:   PopularPosts,
		Archives:       []string{"December 2012 (1)", "December 2011 (2)"},
		Links:          []string{"Check out HM (HEVC reference software)", "Download JCT-VC documents"}}

	return mySideBar, nil
}
