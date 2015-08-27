package controllers

import (
	"html/template"
	"models"
	"net/http"
	"strings"
	"time"
)

func HandlerNew(w http.ResponseWriter, r *http.Request) {
	isAdmin, DisplayName, url, err := models.Login(r)
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
				ArticleID, err := models.FindLatestArticleID(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				a := &models.Article{
					ArticleID:   ArticleID,
					Title:       strings.TrimSpace(r.FormValue("Title")),
					Date:        time.Now(),
					DisplayName: DisplayName,
					Category:    strings.ToLower(r.FormValue("Category")),
					Tags:        strings.Split(strings.ToLower(r.FormValue("Tags")), ","),
					Content:     template.HTML(r.FormValue("Content")),
					Views:       0,
					Comments:    0,
				}

				if err = models.SaveArticle(r, nil, a); err != nil {
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
