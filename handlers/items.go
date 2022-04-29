package handlers

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func GetItem(request *http.Request) (*Item, error) {
	params := mux.Vars(request)
	ItemId, err := strconv.Atoi(params["item_id"])
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow(
		"SELECT * FROM Item WHERE shipment_id IS NULL AND item_id = $1 LIMIT 1;",
		ItemId,
	)

	item := Item{}
	createdAt := time.Time{}
	modifiedAt := time.Time{}
	err = row.Scan(
		&item.ItemId,
		&item.ShipmentId,
		&item.Product,
		&item.Quantity,
		&item.Price,
		&item.Supplier,
		&createdAt,
		&modifiedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func GetItems() ([]Item, error) {
	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Item WHERE shipment_id IS NULL ORDER BY item_id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	createdAt := time.Time{}
	modifiedAt := time.Time{}
	for rows.Next() {
		item := Item{}
		err = rows.Scan(
			&item.ItemId,
			&item.ShipmentId,
			&item.Product,
			&item.Quantity,
			&item.Price,
			&item.Supplier,
			&createdAt,
			&modifiedAt,
		)
		if err != nil {
			return nil, err
		}

		item.CreatedAt = createdAt.Format("2006-01-02")
		item.ModifiedAt = modifiedAt.Format("2006-01-02")
		items = append(items, item)
	}

	return items, nil
}

func PostNewItem(request *http.Request) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer db.Close()

	product := request.FormValue("product")
	quantity := request.FormValue("quantity")
	price := request.FormValue("price")
	supplier := request.FormValue("supplier")

	_, err = db.Query(
		"INSERT INTO Item (product, quantity, price, supplier) VALUES ($1, $2, $3, $4);",
		product,
		quantity,
		price,
		supplier,
	)
	if err != nil {
		return err
	}

	return nil
}

func PostEditItem(request *http.Request) error {
	params := mux.Vars(request)
	itemId, err := strconv.Atoi(params["item_id"])
	if err != nil {
		return err
	}

	err = request.ParseForm()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer db.Close()

	product := request.FormValue("product")
	quantity := request.FormValue("quantity")
	price := request.FormValue("price")
	supplier := request.FormValue("supplier")

	_, err = db.Query(
		"UPDATE Item SET product = $2, quantity = $3, price = $4, supplier = $5, modified_at = NOW() WHERE shipment_id IS NULL AND item_id = $1",
		itemId,
		product,
		quantity,
		price,
		supplier,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteItem(request *http.Request) error {
	params := mux.Vars(request)
	ItemId, err := strconv.Atoi(params["item_id"])
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Query(
		"DELETE FROM Item WHERE shipment_id IS NULL AND item_id = $1;",
		ItemId,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewItemHandler(response http.ResponseWriter, request *http.Request) {
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

func EditItemHandler(response http.ResponseWriter, request *http.Request) {
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

func DeleteItemHandler(response http.ResponseWriter, request *http.Request) {
	var statusCode int

	defer func() {
		LogHTTPTraffic(request, statusCode)
	}()

	switch request.Method {
	case "DELETE":
		err := DeleteItem(request)
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

func ExportCSVHandler(response http.ResponseWriter, request *http.Request) {
	var statusCode int

	defer func() {
		LogHTTPTraffic(request, statusCode)
	}()

	switch request.Method {
	case "GET":
		items, err := GetItems()
		if err != nil {
			fmt.Println(err)
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
			return
		}

		buf := new(bytes.Buffer)
		csvWriter := csv.NewWriter(buf)

		for _, item := range items {
			itemSlice := []string{
				strconv.Itoa(int(item.ItemId)),
				item.Product,
				strconv.Itoa(int(item.Quantity)),
				fmt.Sprintf("%f", item.Price),
				item.Supplier,
				item.CreatedAt,
				item.ModifiedAt,
			}

			err = csvWriter.Write(itemSlice)
			if err != nil {
				fmt.Println(err)
				statusCode = http.StatusInternalServerError
				response.WriteHeader(statusCode)
				return
			}
		}

		csvWriter.Flush()

		response.Header().Set("Content-Disposition", "attachment; filename=items.csv")
		response.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		io.Copy(response, buf)

		statusCode = http.StatusOK
	default:
		statusCode = http.StatusMethodNotAllowed
		response.WriteHeader(statusCode)
	}
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

		err = tmpl.ExecuteTemplate(response, "dashboard.html", items)
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
