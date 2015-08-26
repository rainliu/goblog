package controllers

import (
	"html/template"
	"models"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/user"
)

func HandlerNew(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	if user.IsAdmin(c) {
		if r.Method == "GET" {
			myContent := &models.Content{
				IsAdmin:     true,
				Articles:    nil,
				PrevEntries: "",
				NextEntries: "",
				IsSingle:    false,
				Comments:    nil,
			}

			RenderContent(w, r, myContent, "new")
		} else {
			if r.FormValue("Title") != "" {
				ArticleID, err := models.FindLatestArticleID(c)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				a := &models.Article{
					ArticleID:   ArticleID,
					Title:       strings.TrimSpace(r.FormValue("Title")),
					Date:        time.Now(),
					DisplayName: u.String(),
					Author:      *u,
					Category:    strings.ToLower(r.FormValue("Category")),
					Tags:        strings.Split(strings.ToLower(r.FormValue("Tags")), ","),
					Content:     template.HTML(r.FormValue("Content")),
					Views:       0,
					Comments:    0,
				}

				if err = models.SaveArticle(c, nil, a); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
