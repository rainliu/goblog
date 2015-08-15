package controllers

import (
	"net/http"
	"strconv"

	"models"

	"appengine"
	"appengine/datastore"
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

	var myArticles []models.Article
	var myContent models.Content
	var a *datastore.Query

	r.ParseForm()

	ArticleId, _ := strconv.ParseInt(r.Form["id"][0], 10, 32)
	ArticleID := int(ArticleId)

	a = datastore.NewQuery("Article").Filter("ArticleID =", ArticleID)
	myArticles = make([]models.Article, 0, 1)

	if key, err := a.GetAll(c, &myArticles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		myArticles[0].Views++
		_, err := datastore.Put(c, key[0], &myArticles[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	oldArticleID := int(-1)
	olda := datastore.NewQuery("Article").Filter("ArticleID <", ArticleID).Order("-ArticleID").Limit(1)
	if count, _ := olda.Count(c); count == 1 {
		oldArticles := make([]models.Article, 0, 1)

		if _, err := olda.GetAll(c, &oldArticles); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		oldArticleID = oldArticles[0].ArticleID
	}

	newArticleID := int(-1)
	newa := datastore.NewQuery("Article").Filter("ArticleID >", ArticleID).Order("ArticleID").Limit(1)
	if count, _ := newa.Count(c); count == 1 {
		newArticles := make([]models.Article, 0, 1)

		if _, err := newa.GetAll(c, &newArticles); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newArticleID = newArticles[0].ArticleID
	}

	var myComments []models.Comment
	a = datastore.NewQuery("Comment").Filter("ArticleID =", ArticleID).Order("Date")
	count, _ := a.Count(c)
	myComments = make([]models.Comment, 0, count)

	if _, err := a.GetAll(c, &myComments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//c.Infof(strconv.Itoa(len(myComments)))

	myContent = models.Content{
		IsAdmin:  isAdmin,
		Articles: myArticles,
		PrevEntries: func(id int) string {
			if id < 0 {
				return ""
			} else {
				return "/view/?id=" + strconv.Itoa(id)
			}
		}(oldArticleID),
		NextEntries: func(id int) string {
			if id < 0 {
				return ""
			} else {
				return "/view/?id=" + strconv.Itoa(id)
			}
		}(newArticleID),
		IsSingle: true,
		Comments: myComments,
	}

	RenderContent(w, r, myContent, "content")
}
