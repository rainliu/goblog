package controllers

import (
	"bytes"
	"html/template"
	"models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/user"
)

func RenderContent(w http.ResponseWriter, r *http.Request, myContent *models.Content, tmplName string) {
	var url string

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, _ = user.LoginURL(c, r.URL.String())
	} else {
		url, _ = user.LogoutURL(c, r.URL.String())
	}

	RecentComments, RecentPosts, PopularPosts, err := models.GetSideBarInfo(c)
	if err != nil {
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
	t = t.Funcs(template.FuncMap{"Date2String": Date2String,
		"Author2String":          Author2String,
		"Tags2String":            Tags2String,
		"Comments2String":        Comments2String,
		"CommentContent2Excerpt": CommentContent2Excerpt,
		"FullText2Excerpt":       FullTextSelect(r.Form["id"] != nil),
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

func Date2String(date time.Time) string {
	str := date.Local().String()

	return str[0:10]
}

func Author2String(author user.User) string {
	return author.String()
}

func Tags2String(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	var strBuffer bytes.Buffer
	for i := 0; i < len(tags); i++ {
		strBuffer.WriteString(strings.TrimSpace(tags[i]))
		if i != len(tags)-1 {
			strBuffer.WriteString(", ")
		}
	}

	return strBuffer.String()
}

func FullTextSelect(single bool) func(template.HTML) template.HTML {
	if single {
		return FullText2Full
	} else {
		return FullText2Excerpt
	}
}

func FullText2Full(text template.HTML) template.HTML {
	return text
}

func FullText2Excerpt(text template.HTML) template.HTML {
	if len(text) < models.MyBlogConfig.ExcerptsCharNumber {
		return text
	} else {
		return text[0:models.MyBlogConfig.ExcerptsCharNumber] + "[...]"
	}
}

func Comments2String(comments []string) string {
	return strconv.Itoa(len(comments))
}

func CommentContent2Excerpt(text string) string {
	if len(text) < models.MyBlogConfig.ExcerptsCommentsCharNumber {
		return text
	} else {
		return text[0:models.MyBlogConfig.ExcerptsCommentsCharNumber] + " [...]"
	}
}
