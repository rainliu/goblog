package controllers

import (
	"html/template"
	"models"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func RenderContent(w http.ResponseWriter, r *http.Request, myContent models.Content, tmplName string) {
	var url string
	var err error

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, _ = user.LoginURL(c, r.URL.String())
	} else {
		url, _ = user.LogoutURL(c, r.URL.String())
	}

	a := datastore.NewQuery("Comment").Order("-Date").Limit(models.MyBlogConfig.ExcerptsNumber)
	RecentComments := make([]models.Comment, 0, models.MyBlogConfig.ExcerptsNumber)
	if _, err := a.GetAll(c, &RecentComments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := datastore.NewQuery("Article").Order("-Date").Limit(models.MyBlogConfig.ExcerptsNumber)
	RecentPosts := make([]models.Article, 0, models.MyBlogConfig.ExcerptsNumber)
	if _, err := b.GetAll(c, &RecentPosts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var PopularPosts []models.Article
	p := datastore.NewQuery("Article").Order("-Views").Limit(models.MyBlogConfig.ExcerptsNumber)
	PopularPosts = make([]models.Article, 0, models.MyBlogConfig.ExcerptsNumber)
	if _, err := p.GetAll(c, &PopularPosts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mySideBar := models.SideBar{
		User:           u,
		LoginoutURL:    url,
		RSSConfig:      models.MyBlogConfig,
		RecentComments: RecentComments,
		RecentPosts:    RecentPosts,
		PopularPosts:   PopularPosts,
		Archives:       []string{"December 2012 (1)", "December 2011 (2)"},
		Links:          []string{"Check out HM (HEVC reference software)", "Download JCT-VC documents"}}

	w.Header().Set("Content-Type", "text/html")

	t := template.New("home")
	t = t.Funcs(template.FuncMap{"Date2String": models.Date2String,
		"Author2String":          models.Author2String,
		"Tags2String":            models.Tags2String,
		"Comments2String":        models.Comments2String,
		"CommentContent2Excerpt": models.CommentContent2Excerpt,
		"FullText2Excerpt":       models.FullTextSelect(r.Form["id"] != nil), //models.FullText2Excerpt,
	})

	t, err = t.ParseFiles("views/header.tmpl", "views/"+tmplName+".tmpl", "views/sidebar.tmpl", "views/footer.tmpl")

	err = t.ExecuteTemplate(w, "header", models.MyBlogConfig)
	err = t.ExecuteTemplate(w, tmplName, myContent)
	err = t.ExecuteTemplate(w, "sidebar", mySideBar)
	err = t.ExecuteTemplate(w, "footer", models.MyBlogConfig)

	//err = t.Execute(w, nil)
	if err != nil {
		c.Errorf("%v", err)
	}
}
