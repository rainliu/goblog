package controllers

import (
	"models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
)

func HandlerComment(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		if strings.TrimSpace(r.FormValue("Name")) != "" && strings.TrimSpace(r.FormValue("Email")) != "" {
			ArticleId, _ := strconv.ParseInt(r.FormValue("ArticleID"), 10, 32)
			ArticleID := int(ArticleId)

			err := models.UpdateArticleComments(c, ArticleID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			CommentID, err := models.FindLatestCommentID(c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			comment := &models.Comment{
				CommentID: CommentID,
				ArticleID: ArticleID,
				Date:      time.Now(),
				Name:      strings.TrimSpace(r.FormValue("Name")),
				Email:     strings.TrimSpace(r.FormValue("Email")),
				Website:   strings.TrimSpace(r.FormValue("Website")),
				Content:   r.FormValue("Content"),
			}

			err = models.SaveComment(c, nil, comment)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
