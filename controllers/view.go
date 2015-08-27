package controllers

import (
	"net/http"
	"strconv"

	"models"
)

func HandlerView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ArticleId, _ := strconv.ParseInt(r.Form["id"][0], 10, 32)
	ArticleID := int(ArticleId)

	myContent, err := models.GetContent(r, ArticleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderContent(w, r, myContent, "content")
}
