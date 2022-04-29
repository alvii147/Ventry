\connect ventrydb;

CREATE TABLE Shipment (
    shipment_id INT GENERATED ALWAYS AS IDENTITY,
    shipper VARCHAR(50) NOT NULL,
    receiver VARCHAR(50) NOT NULL,
    shipped_at DATE NOT NULL DEFAULT CURRENT_DATE,
    delivered_at DATE NOT NULL,

    PRIMARY KEY (shipment_id)
);

CREATE TABLE Item (
    item_id INT GENERATED ALWAYS AS IDENTITY,
    shipment_id INT,
    product VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    price FLOAT8 NOT NULL,
    supplier VARCHAR(50) NOT NULL,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    modified_at DATE NOT NULL DEFAULT CURRENT_DATE,

    PRIMARY KEY (item_id),
    CONSTRAINT fk_shipment FOREIGN KEY (shipment_id) REFERENCES Shipment(shipment_id)
);
