\connect ventrydb;

CREATE TABLE Inventory (
    item_id SERIAL PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    price FLOAT8 NOT NULL,
    owner VARCHAR(50) NOT NULL,
    supplier VARCHAR(50) NOT NULL,
    shipper VARCHAR(50),
    created DATE NOT NULL DEFAULT CURRENT_DATE,
    modified DATE NOT NULL DEFAULT CURRENT_DATE
);
