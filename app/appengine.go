package app

import (
	"controllers"
	"io"
	"net/http"

	"appengine"
)

func init() {
	http.HandleFunc("/", controllers.HandlerHome)
	http.HandleFunc("/new", controllers.HandlerNew)
	http.HandleFunc("/delete/", controllers.HandlerDelete)
	http.HandleFunc("/edit/", controllers.HandlerEdit)
	http.HandleFunc("/view/", controllers.HandlerView)
	http.HandleFunc("/comment", controllers.HandlerComment)
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
}
