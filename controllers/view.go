package controllers

import (
	"net/http"
	"strconv"

	"models"

	"appengine"
	"appengine/user"
)

func HandlerView(w http.ResponseWriter, r *http.Request) {
	var isAdmin bool

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		isAdmin = false
	} else {
		isAdmin = user.IsAdmin(c)
	}

	r.ParseForm()
	ArticleId, _ := strconv.ParseInt(r.Form["id"][0], 10, 32)
	ArticleID := int(ArticleId)

	myContent, err := models.GetContent(c, isAdmin, ArticleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderContent(w, r, myContent, "content")
}
