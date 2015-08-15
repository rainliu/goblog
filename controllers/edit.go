package controllers

import (
	"html/template"
	"models"
	"net/http"
	"strconv"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func HandlerEdit(w http.ResponseWriter, r *http.Request) {
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
		var err error
		var key []*datastore.Key
		var myArticles []models.Article
		r.ParseForm()
		if r.Form["id"] != nil {
			ArticleId, _ := strconv.ParseInt(r.Form["id"][0], 10, 32)
			ArticleID := int(ArticleId)

			a := datastore.NewQuery("Article").Filter("ArticleID =", ArticleID)
			myArticles = make([]models.Article, 0, 1)

			if key, err = a.GetAll(c, &myArticles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if r.Method == "GET" {
				myContent := models.Content{
					IsAdmin:     true,
					Articles:    myArticles,
					PrevEntries: "",
					NextEntries: "",
					IsSingle:    false,
					Comments:    nil,
				}

				RenderContent(w, r, myContent, "edit")
			} else {
				ArticleId, _ := strconv.ParseInt(r.FormValue("ArticleID"), 10, 32)
				ArticleID := int(ArticleId)
				if myArticles[0].ArticleID != ArticleID {
					http.Error(w, "Wrong ArticleID", http.StatusInternalServerError)
					return
				}

				a := models.Article{
					ArticleID:   ArticleID,
					Title:       strings.TrimSpace(r.FormValue("Title")),
					Date:        myArticles[0].Date,
					DisplayName: u.String(),
					Author:      *u,
					Category:    strings.ToLower(r.FormValue("Category")),
					Tags:        strings.Split(strings.ToLower(r.FormValue("Tags")), ","),
					Content:     template.HTML(r.FormValue("Content")),
					Views:       myArticles[0].Views,
					Comments:    myArticles[0].Comments,
				}
				_, err := datastore.Put(c, key[0], &a)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, "/", http.StatusFound)
			}
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
