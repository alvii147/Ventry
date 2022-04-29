package handlers

import (
	"fmt"
	"net/http"
)

type dashboardContext struct {
	Items     []Item
	Shipments []Shipment
}

func DashboardHandler(response http.ResponseWriter, request *http.Request) {
	var statusCode int

	defer func() {
		LogHTTPTraffic(request, statusCode)
	}()

	switch request.Method {
	case "GET":
		statusCode = http.StatusOK

		items, err := GetItems()
		if err != nil {
			fmt.Println(err)
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
			return
		}

		shipments, err := GetShipments()
		if err != nil {
			fmt.Println(err)
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
			return
		}

		ctx := dashboardContext{
			items,
			shipments,
		}

		err = tmpl.ExecuteTemplate(response, "dashboard.html", ctx)
		if err != nil {
			fmt.Println(err)
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
			return
		}
	default:
		statusCode = http.StatusMethodNotAllowed
		response.WriteHeader(statusCode)
	}
}
