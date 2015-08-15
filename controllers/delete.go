package controllers

import (
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	if user.IsAdmin(c) {
		r.ParseForm()
		if r.Form["id"] != nil {
			ArticleId, _ := strconv.ParseInt(r.Form["id"][0], 10, 32)
			ArticleID := int(ArticleId)

			a := datastore.NewQuery("Article").Filter("ArticleID =", ArticleID).KeysOnly()

			if k, err := a.GetAll(c, nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				if err = datastore.Delete(c, k[0]); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
