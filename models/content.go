package models

import (
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Content struct {
	IsAdmin     bool
	Articles    []Article
	PrevEntries string
	NextEntries string
	IsSingle    bool
	Comments    []Comment
}

func GetContents(r *http.Request, PageID int) (*Content, error) {
	var isAdmin bool

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		isAdmin = false
	} else {
		isAdmin = user.IsAdmin(c)
	}

	TotalArticles, _ := datastore.NewQuery("Article").Count(c)

	a := datastore.NewQuery("Article").Order("-Date").Offset(PageID * MyBlogConfig.ExcerptsNumber).Limit(MyBlogConfig.ExcerptsNumber)

	myArticles := make([]Article, 0, MyBlogConfig.ExcerptsNumber)

	if _, err := a.GetAll(c, &myArticles); err != nil {
		return nil, err
	}

	myContent := &Content{
		IsAdmin:  isAdmin,
		Articles: myArticles,
		PrevEntries: func(id int) string {
			if (id+1)*MyBlogConfig.ExcerptsNumber >= TotalArticles {
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
	return myContent, nil
}

func GetContent(r *http.Request, ArticleID int) (*Content, error) {
	var isAdmin bool

	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		isAdmin = false
	} else {
		isAdmin = user.IsAdmin(c)
	}

	a := datastore.NewQuery("Article").Filter("ArticleID =", ArticleID)
	myArticles := make([]Article, 0, 1)

	if key, err := a.GetAll(c, &myArticles); err != nil {
		return nil, err
	} else {
		myArticles[0].Views++
		_, err := datastore.Put(c, key[0], &myArticles[0])
		if err != nil {
			return nil, err
		}
	}

	oldArticleID := int(-1)
	olda := datastore.NewQuery("Article").Filter("ArticleID <", ArticleID).Order("-ArticleID").Limit(1)
	if count, _ := olda.Count(c); count == 1 {
		oldArticles := make([]Article, 0, 1)

		if _, err := olda.GetAll(c, &oldArticles); err != nil {
			return nil, err
		}
		oldArticleID = oldArticles[0].ArticleID
	}

	newArticleID := int(-1)
	newa := datastore.NewQuery("Article").Filter("ArticleID >", ArticleID).Order("ArticleID").Limit(1)
	if count, _ := newa.Count(c); count == 1 {
		newArticles := make([]Article, 0, 1)

		if _, err := newa.GetAll(c, &newArticles); err != nil {
			return nil, err
		}
		newArticleID = newArticles[0].ArticleID
	}

	a = datastore.NewQuery("Comment").Filter("ArticleID =", ArticleID).Order("Date")
	count, _ := a.Count(c)
	myComments := make([]Comment, 0, count)

	if _, err := a.GetAll(c, &myComments); err != nil {
		return nil, err
	}

	myContent := &Content{
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

	return myContent, nil
}
