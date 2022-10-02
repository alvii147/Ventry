<p align="center">
    <img alt="Ventry Logo" src="static/images/ventry_logo.png" width="300" />
</p>

<p align="center">
    <a href="https://goreportcard.com/report/github.com/alvii147/Ventry" target="_blank" rel="noopener noreferrer" style="text-decoration: none;">
        <img alt="Go Report Card Badge" src="https://goreportcard.com/badge/github.com/alvii147/Ventry" />
    </a>
    <a href="https://ventry.zahinzaman1.repl.co/" target="_blank" rel="noopener noreferrer" style="text-decoration: none;">
        <img alt="Live Demo" src="https://img.shields.io/badge/replit-Live%20Demo-steelblue?logo=replit" />
    </a>
</p>

# Overview

Ventry is an inventory tracking [CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete) application built using [Go](https://go.dev/) and [PostgreSQL](https://www.postgresql.org/). This is a submission for the **Shopify Fall 2022 Backend Developer Internship Challenge**.

Check out the hosted application [here!](https://ventry.zahinzaman1.repl.co/)

## Tech Stack

![Tech Stack](img/techstack_bg.png)

## Dashboard

The dashboard page shows all items in the current inventory as well as all currently scheduled shipments.

![Dashboard Screenshot](img/dashboard_screenshot.png)

## Create New Item

New items can be created by navigating to the **New Item** link on the navbar and filling out the form with the item details.

<img alt="New Item Screenshot" src="img/new_item_screenshot.png" width="500" />

## Update Existing Item

Existing items can be updated by clicking the edit icon on the right hand side of the corresponding item on the dashboard and filling out the pre-populated form with the appropriate changes.

<img alt="Edit Item Screenshot" src="img/edit_item_screenshot.png" width="500" />

## Create New Shipment

New shipments can be created by navigating to the **New Shipment** link on the navbar and filling out the form with the shipping details, including which items to ship.

<img alt="New Shipment Screenshot" src="img/new_shipment_screenshot.png" width="500" />

## Exporting to CSV

Ventry can also export the inventory data into [CSV](https://en.wikipedia.org/wiki/Comma-separated_values) format. This can be done by clicking the **Export CSV** button on the dashboard.

<img alt="Export CSV Screenshot" src="img/export_screenshot.png" width="350" />

## Database Schema

<img alt="Database Schema" src="img/schema.png" width="600" />

# Getting Started

## Running Ventry using Docker Compose (Recommended)

:one: Install Docker from the [official website](https://www.docker.com/). Installation instructions may vary depending on the OS.

:two: Run Docker.

```bash
docker compose up
```

Once the logs print `database system is ready to accept connections`, Ventry will be up on `http://localhost:8000/`.

## Running Ventry using Go & PostgreSQL

:one: Install Go from the [official website](https://go.dev/). Installation instructions may vary depending on the OS.

:two: Install PostgreSQL from the [official website](https://www.postgresql.org/). Installation instructions may vary depending on the OS.

:three: Set environment variables.

```bash
export VENTRY_POSTGRES_DATABASE="ventrydb"
export VENTRY_POSTGRES_USERNAME="<postgres username>"
export VENTRY_POSTGRES_PASSWORD="<postgres password>"
export VENTRY_POSTGRES_HOST="localhost"
export VENTRY_POSTGRES_PORT="5432"
```

Note that `VENTRY_POSTGRES_USERNAME` and `VENTRY_POSTGRES_PASSWORD` must correspond to the username and passwords of the PostgreSQL instance running respectively.

:four: Set up database.

Run the SQL query scripts in `db/` to create the database and table.

```bash
psql -U <postgres username> -d ventrydb -a -f db/create_database.sql
psql -U <postgres username> -d ventrydb -a -f db/create_tables.sql
psql -U <postgres username> -d ventrydb -a -f db/populate_tables.sql
```

:five: Install Go dependencies.

```bash
go mod download
```

:six: Run Go server.

```bash
go run .
```

This should run the Ventry web app on `http://localhost:8000/`.
