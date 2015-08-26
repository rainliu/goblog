package controllers

import (
	"net/http"
	"strconv"

	"models"

	"appengine"
	"appengine/user"
)

func HandlerHome(w http.ResponseWriter, r *http.Request) {
	var isAdmin bool

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		isAdmin = false
	} else {
		isAdmin = user.IsAdmin(c)
	}

	PageID := int(0)

	if r.Form["page"] != nil {
		PageId, _ := strconv.ParseInt(r.Form["page"][0], 10, 32)
		PageID = int(PageId)
	}

	myContent, err := models.GetContents(c, isAdmin, PageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderContent(w, r, myContent, "content")
}
