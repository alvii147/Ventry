package handlers

import (
	"fmt"
	"net/http"
)

func EditHandler(response http.ResponseWriter, request *http.Request) {
	var statusCode int

	defer func() {
		LogHTTPTraffic(request, statusCode)
	}()

	switch request.Method {
	case "GET":
		statusCode = http.StatusOK

		item, err := GetItem(request)
		if err != nil {
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
		}

		err = tmpl.ExecuteTemplate(response, "edit.html", item)
		if err != nil {
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
		}
	case "POST":
		err := PostEditItem(request)
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
