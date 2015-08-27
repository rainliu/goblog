package app

import (
	"controllers"
	"net/http"
)

func init() {
	http.HandleFunc("/", controllers.HandlerHome)
	http.HandleFunc("/new", controllers.HandlerNew)
	http.HandleFunc("/delete/", controllers.HandlerDelete)
	http.HandleFunc("/edit/", controllers.HandlerEdit)
	http.HandleFunc("/view/", controllers.HandlerView)
	http.HandleFunc("/comment", controllers.HandlerComment)
}
