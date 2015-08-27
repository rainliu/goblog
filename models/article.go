package models

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

type Article struct {
	ArticleID   int
	Title       string
	Date        time.Time
	DisplayName string
	Category    string
	Tags        []string
	Content     template.HTML `datastore:"Content,noindex"`
	Views       int
	Comments    int
}

func FindLatestArticleID(r *http.Request) (int, error) {
	c := appengine.NewContext(r)

	ArticleID := 0
	last := datastore.NewQuery("Article")
	if count, err := last.Count(c); count == 0 || err != nil {
		ArticleID = 0
	} else {
		last = last.Order("-Date").Limit(1)
		lastArticle := make([]Article, 0, 1)
		if _, err := last.GetAll(c, &lastArticle); err != nil {
			return 0, err
		}
		ArticleID = lastArticle[0].ArticleID + 1
	}
	return ArticleID, nil
}

func FindArticle(r *http.Request, ArticleID int) ([]Article, []*datastore.Key, error) {
	c := appengine.NewContext(r)

	a := datastore.NewQuery("Article").Filter("ArticleID =", ArticleID)
	myArticles := make([]Article, 0, 1)

	if keys, err := a.GetAll(c, &myArticles); err != nil {
		return nil, nil, err
	} else {
		return myArticles, keys, nil
	}
}

func SaveArticle(r *http.Request, key *datastore.Key, a *Article) error {
	c := appengine.NewContext(r)

	if key == nil {
		key = datastore.NewIncompleteKey(c, "Article", nil)
	}
	_, err := datastore.Put(c, key, a)
	return err
}

func DeleteArticle(r *http.Request, ArticleID int) error {
	c := appengine.NewContext(r)

	a := datastore.NewQuery("Article").Filter("ArticleID =", ArticleID).KeysOnly()

	if k, err := a.GetAll(c, nil); err != nil {
		return err
	} else {
		if err = datastore.Delete(c, k[0]); err != nil {
			return err
		}
	}
	return nil
}

func UpdateArticleComments(r *http.Request, ArticleID int) error {
	c := appengine.NewContext(r)

	a := datastore.NewQuery("Article").Filter("ArticleID =", ArticleID)
	myArticles := make([]Article, 0, 1)

	if key, err := a.GetAll(c, &myArticles); err != nil {
		return err
	} else {
		myArticles[0].Comments++
		_, err := datastore.Put(c, key[0], &myArticles[0])
		if err != nil {
			return err
		}
	}
	return nil
}
