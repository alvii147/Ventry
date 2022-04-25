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
		"SELECT * FROM Inventory WHERE item_id = $1 LIMIT 1;",
		ItemId,
	)

	item := Item{}
	created := time.Time{}
	modified := time.Time{}
	err = row.Scan(
		&item.ItemId,
		&item.Title,
		&item.Quantity,
		&item.Price,
		&item.Owner,
		&item.Supplier,
		&item.Shipper,
		&created,
		&modified,
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

	rows, err := db.Query("SELECT * FROM Inventory ORDER BY item_id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Item, 0)
	created := time.Time{}
	modified := time.Time{}
	for rows.Next() {
		item := Item{}
		err = rows.Scan(
			&item.ItemId,
			&item.Title,
			&item.Quantity,
			&item.Price,
			&item.Owner,
			&item.Supplier,
			&item.Shipper,
			&created,
			&modified,
		)
		if err != nil {
			return nil, err
		}

		item.Created = created.Format("2006-01-02")
		item.Modified = modified.Format("2006-01-02")
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

	title := request.FormValue("title")
	quantity := request.FormValue("quantity")
	price := request.FormValue("price")
	owner := request.FormValue("owner")
	supplier := request.FormValue("supplier")
	shipper := request.FormValue("shipper")

	_, err = db.Query(
		"INSERT INTO Inventory (title, quantity, price, owner, supplier, shipper) VALUES ($1, $2, $3, $4, $5, $6)",
		title,
		quantity,
		price,
		owner,
		supplier,
		shipper,
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

	title := request.FormValue("title")
	quantity := request.FormValue("quantity")
	price := request.FormValue("price")
	owner := request.FormValue("owner")
	supplier := request.FormValue("supplier")
	shipper := request.FormValue("shipper")

	_, err = db.Query(
		"UPDATE Inventory SET title = $2, quantity = $3, price = $4, owner = $5, supplier = $6, shipper = $7, modified = NOW() WHERE item_id = $1",
		itemId,
		title,
		quantity,
		price,
		owner,
		supplier,
		shipper,
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
		"DELETE FROM Inventory WHERE item_id = $1;",
		ItemId,
	)
	if err != nil {
		return err
	}

	return nil
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
				item.Title,
				strconv.Itoa(int(item.Quantity)),
				fmt.Sprintf("%f", item.Price),
				item.Owner,
				item.Supplier,
				item.Shipper,
				item.Created,
				item.Modified,
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

		response.Header().Set("Content-Disposition", "attachment; filename=inventory.csv")
		response.Header().Set("Content-Type", request.Header.Get("Content-Type"))
		io.Copy(response, buf)

		statusCode = http.StatusOK
	default:
		statusCode = http.StatusMethodNotAllowed
		response.WriteHeader(statusCode)
	}
}

func APIHandler(response http.ResponseWriter, request *http.Request) {
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
