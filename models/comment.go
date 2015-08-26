package models

import (
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

func FindLatestCommentID(c appengine.Context) (int, error) {
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

func SaveComment(c appengine.Context, key *datastore.Key, comment *Comment) error {
	if key == nil {
		key = datastore.NewIncompleteKey(c, "Comment", nil)
	}
	_, err := datastore.Put(c, key, comment)
	if err != nil {
		return err
	}
	return nil
}
