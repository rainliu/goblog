package controllers

import (
	"models"
	"net/http"
	"strconv"
)

func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	isAdmin, _, url, err := models.Login(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if url != "" {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	if isAdmin {
		r.ParseForm()
		if r.Form["id"] != nil {
			ArticleId, _ := strconv.ParseInt(r.Form["id"][0], 10, 32)
			ArticleID := int(ArticleId)

			if err := models.DeleteComment(r, ArticleID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := models.DeleteArticle(r, ArticleID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
