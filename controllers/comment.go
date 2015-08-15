package controllers

import (
	"models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
)

func HandlerComment(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		if strings.TrimSpace(r.FormValue("Name")) != "" && strings.TrimSpace(r.FormValue("Email")) != "" {
			ArticleId, _ := strconv.ParseInt(r.FormValue("ArticleID"), 10, 32)
			ArticleID := int(ArticleId)

			a := datastore.NewQuery("Article").Filter("ArticleID =", ArticleID)
			myArticles := make([]models.Article, 0, 1)

			if key, err := a.GetAll(c, &myArticles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				myArticles[0].Comments++
				_, err := datastore.Put(c, key[0], &myArticles[0])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			CommentID := 0
			last := datastore.NewQuery("Comment")
			if count, err := last.Count(c); count == 0 || err != nil {
				CommentID = 0
			} else {
				last = last.Order("-Date").Limit(1)
				lastComment := make([]models.Comment, 0, 1)
				if _, err := last.GetAll(c, &lastComment); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				CommentID = lastComment[0].CommentID + 1
			}

			comment := models.Comment{
				CommentID: CommentID,
				ArticleID: ArticleID,
				Date:      time.Now(),
				Name:      strings.TrimSpace(r.FormValue("Name")),
				Email:     strings.TrimSpace(r.FormValue("Email")),
				Website:   strings.TrimSpace(r.FormValue("Website")),
				Content:   r.FormValue("Content"),
			}

			key := datastore.NewIncompleteKey(c, "Comment", nil)
			_, err := datastore.Put(c, key, &comment)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
