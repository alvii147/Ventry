package handlers

import (
	"fmt"
	"net/http"
)

func NewHandler(response http.ResponseWriter, request *http.Request) {
	var statusCode int

	defer func() {
		LogHTTPTraffic(request, statusCode)
	}()

	switch request.Method {
	case "GET":
		statusCode = http.StatusOK

		err := tmpl.ExecuteTemplate(response, "new.html", nil)
		if err != nil {
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
		}
	case "POST":
		err := PostNewItem(request)
		if err != nil {
			fmt.Println(err)
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
			return
		}

		statusCode = http.StatusSeeOther
		http.Redirect(response, request, "/", statusCode)
	default:
		statusCode = http.StatusMethodNotAllowed
		response.WriteHeader(statusCode)
	}
}
