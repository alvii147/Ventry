package main

import (
	"net/http"

	"github.com/alvii147/Ventry/handlers"
	"github.com/gorilla/mux"
)

const STATIC_FILES_PATH string = "static"

func routing() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/" + STATIC_FILES_PATH + "/").Handler(http.StripPrefix("/"+STATIC_FILES_PATH+"/", http.FileServer(http.Dir(STATIC_FILES_PATH))))
	router.PathPrefix("/favicon.ico").Handler(http.NotFoundHandler())
	router.HandleFunc("/items/new", handlers.NewItemHandler)
	router.HandleFunc("/items/edit/{item_id:[0-9]+}", handlers.EditItemHandler)
	router.HandleFunc("/items/delete/{item_id:[0-9]+}", handlers.DeleteItemHandler)
	router.HandleFunc("/items/export", handlers.ExportCSVHandler)
	router.HandleFunc("/", handlers.DashboardHandler)
	router.HandleFunc("/shipments/new", handlers.NewShipmentHandler)

	return router
}

func main() {
	router := routing()
	http.ListenAndServe(":8000", router)
}
