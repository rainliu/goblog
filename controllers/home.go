package controllers

import (
	"net/http"
	"strconv"

	"models"

	"appengine"
	"appengine/datastore"
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

	var myArticles []models.Article
	var myContent models.Content
	var a *datastore.Query
	TotalArticles, _ := datastore.NewQuery("Article").Count(c)

	PageID := int(0)

	if r.Form["page"] != nil {
		PageId, _ := strconv.ParseInt(r.Form["page"][0], 10, 32)
		PageID = int(PageId)
	}
	a = datastore.NewQuery("Article").Ancestor(newblogKey(c)).Order("-Date").Offset(PageID * models.MyBlogConfig.ExcerptsNumber).Limit(models.MyBlogConfig.ExcerptsNumber)

	myArticles = make([]models.Article, 0, models.MyBlogConfig.ExcerptsNumber)

	if _, err := a.GetAll(c, &myArticles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	myContent = models.Content{
		IsAdmin:  isAdmin,
		Articles: myArticles,
		PrevEntries: func(id int) string {
			if (id+1)*models.MyBlogConfig.ExcerptsNumber >= TotalArticles {
				return ""
			} else {
				return "/?page=" + strconv.Itoa(id+1)
			}
		}(PageID),
		NextEntries: func(id int) string {
			if id == 0 {
				return ""
			} else {
				return "/?page=" + strconv.Itoa(id-1)
			}
		}(PageID),
		IsSingle: false,
		Comments: nil,
	}

	RenderContent(w, r, myContent)
}
