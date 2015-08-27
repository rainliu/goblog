package controllers

import (
	"html/template"
	"models"
	"net/http"
	"strconv"
	"strings"
)

func HandlerEdit(w http.ResponseWriter, r *http.Request) {
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
		r.ParseForm()
		if r.Form["id"] != nil {
			ArticleId, _ := strconv.ParseInt(r.Form["id"][0], 10, 32)
			ArticleID := int(ArticleId)

			myArticles, keys, err := models.FindArticle(r, ArticleID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if r.Method == "GET" {
				myContent := &models.Content{
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

				a := &models.Article{
					ArticleID:   ArticleID,
					Title:       strings.TrimSpace(r.FormValue("Title")),
					Date:        myArticles[0].Date,
					DisplayName: DisplayName,
					Category:    strings.ToLower(r.FormValue("Category")),
					Tags:        strings.Split(strings.ToLower(r.FormValue("Tags")), ","),
					Content:     template.HTML(r.FormValue("Content")),
					Views:       myArticles[0].Views,
					Comments:    myArticles[0].Comments,
				}

				if err = models.SaveArticle(r, keys[0], a); err != nil {
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
