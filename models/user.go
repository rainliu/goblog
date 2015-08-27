package models

import (
	"net/http"

	"appengine"
	"appengine/user"
)

func Login(r *http.Request) (bool, string, string, error) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			return false, "", "", err
		}
		return false, "", url, nil
	}

	return user.IsAdmin(c), u.String(), "", nil
}
