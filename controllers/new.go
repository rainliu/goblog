package controllers

import (
	"html/template"
	"models"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
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
			myContent := models.Content{
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
				ArticleID := 0
				last := datastore.NewQuery("Article")
				if count, err := last.Count(c); count == 0 || err != nil {
					ArticleID = 0
				} else {
					last = last.Ancestor(newblogKey(c)).Order("-Date").Limit(1)
					lastArticle := make([]models.Article, 0, 1)
					if _, err := last.GetAll(c, &lastArticle); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					ArticleID = lastArticle[0].ArticleID + 1
				}

				a := models.Article{
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

				// We set the same parent key on every Greeting entity to ensure each Greeting
				// is in the same entity group. Queries across the single entity group
				// will be consistent. However, the write rate to a single entity group
				// should be limited to ~1/second.
				key := datastore.NewIncompleteKey(c, "Article", newblogKey(c))
				_, err := datastore.Put(c, key, &a)
				if err != nil {
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

// newblogKey returns the key used for all blog entries.
func newblogKey(c appengine.Context) *datastore.Key {
	// The string "default_blog" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Article", "default_article", 0, nil)
}
