package controllers

import (
	"net/http"
	"strconv"

	"models"
)

func HandlerHome(w http.ResponseWriter, r *http.Request) {
	PageID := int(0)

	if r.Form["page"] != nil {
		PageId, _ := strconv.ParseInt(r.Form["page"][0], 10, 32)
		PageID = int(PageId)
	}

	myContent, err := models.GetContents(r, PageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderContent(w, r, myContent, "content")
}
