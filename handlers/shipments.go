package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func GetShipmentItems(itemId int) ([]Item, error) {
	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		"SELECT item_id, shipment_id, product, quantity, price, supplier, created_at, modified_at FROM Item WHERE shipment_id = $1 ORDER BY item_id;",
		itemId,
	)
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

func GetShipments() ([]Shipment, error) {
	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT shipment_id, shipper, receiver, shipped_at, delivered_at FROM Shipment ORDER BY shipment_id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shipments := make([]Shipment, 0)
	shippedAt := time.Time{}
	deliveredAt := time.Time{}
	for rows.Next() {
		shipment := Shipment{}
		err = rows.Scan(
			&shipment.ShipmentId,
			&shipment.Shipper,
			&shipment.Receiver,
			&shippedAt,
			&deliveredAt,
		)
		if err != nil {
			return nil, err
		}

		items, err := GetShipmentItems(shipment.ShipmentId)
		if err != nil {
			return nil, err
		}
		shipment.Items = items

		shipment.ShippedAt = shippedAt.Format("2006-01-02")
		shipment.DeliveredAt = deliveredAt.Format("2006-01-02")
		daysLeft := int(time.Until(deliveredAt))
		if daysLeft < 0 {
			daysLeft = 0
		}
		shipment.DaysLeft = daysLeft
		shipment.DaysTotal = int(deliveredAt.Sub(shippedAt).Hours() / 24)

		shipments = append(shipments, shipment)
	}

	return shipments, nil
}

func PostNewShipment(request *http.Request) error {
	err := request.ParseForm()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", POSTGRES_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer db.Close()

	shipper := request.FormValue("shipper")
	receiver := request.FormValue("receiver")
	deliveredAt := request.FormValue("delivered_at")

	fmt.Println(shipper, receiver, deliveredAt)
	row := db.QueryRow(
		"INSERT INTO Shipment (shipper, receiver, delivered_at) VALUES ($1, $2, $3) RETURNING shipment_id;",
		shipper,
		receiver,
		deliveredAt,
	)
	var shipmentId int
	err = row.Scan(&shipmentId)
	fmt.Println(shipmentId)

	if err != nil {
		return err
	}

	re := regexp.MustCompile(`^item-(\d+)$`)
	for key := range request.Form {
		matches := re.FindStringSubmatch(key)
		if len(matches) < 2 {
			continue
		}

		itemId, err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}

		fmt.Println(itemId)
		_, err = db.Query(
			"UPDATE Item SET shipment_id = $1 WHERE shipment_id IS NULL AND item_id = $2",
			shipmentId,
			itemId,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewShipmentHandler(response http.ResponseWriter, request *http.Request) {
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

		err = tmpl.ExecuteTemplate(response, "shipment_new.html", items)
		if err != nil {
			fmt.Println(err)
			statusCode = http.StatusInternalServerError
			response.WriteHeader(statusCode)
			return
		}
	case "POST":
		err := PostNewShipment(request)
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
