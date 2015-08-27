package models

import (
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

type Comment struct {
	CommentID int
	ArticleID int
	Date      time.Time
	Name      string
	Email     string
	Website   string
	Content   string
}

func FindLatestCommentID(r *http.Request) (int, error) {
	c := appengine.NewContext(r)

	CommentID := 0
	last := datastore.NewQuery("Comment")
	if count, err := last.Count(c); count == 0 || err != nil {
		CommentID = 0
	} else {
		last = last.Order("-Date").Limit(1)
		lastComment := make([]Comment, 0, 1)
		if _, err := last.GetAll(c, &lastComment); err != nil {
			return 0, err
		}
		CommentID = lastComment[0].CommentID + 1
	}

	return CommentID, nil
}

func DeleteComment(r *http.Request, ArticleID int) error {
	c := appengine.NewContext(r)

	a := datastore.NewQuery("Comment").Filter("ArticleID =", ArticleID).KeysOnly()

	if k, err := a.GetAll(c, nil); err != nil {
		return err
	} else {
		if err = datastore.Delete(c, k[0]); err != nil {
			return err
		}
	}
	return nil
}

func SaveComment(r *http.Request, comment *Comment) error {
	c := appengine.NewContext(r)

	key := datastore.NewIncompleteKey(c, "Comment", nil)

	_, err := datastore.Put(c, key, comment)
	if err != nil {
		return err
	}
	return nil
}
