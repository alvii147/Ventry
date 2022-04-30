\connect ventrydb;

INSERT INTO Shipment
    (shipment_id, shipper, receiver, shipped_at, delivered_at)
VALUES
    (1, 'Instacart', 'Walmart', '2022-04-20', '2022-05-28'),
    (2, 'FedEx', 'No Frills', '2022-04-26', '2022-05-17');

INSERT INTO Item
    (item_id, shipment_id, product, quantity, price, supplier)
VALUES
    (1, NULL, 'Cream Cheese', 75, 3.69, 'PHILADELPHIA'),
    (2, NULL, 'Eggs', 144, 4.29, 'Poultry Farm'),
    (3, NULL, 'Yogurt', 20, 5.52, 'Astro'),
    (4, 1, 'Dates', 50, 0.12, 'Karamah-Jordan Valley Farm'),
    (5, 2, 'Potatoes', 68, 0.38, 'Local Farm'),
    (6, 1, 'Tea bags', 108, 0.14, 'Sylhet Tea Garden'),
    (7, 1, 'Hummus', 68, 0.38, 'Kabul Farms'),
    (8, 2, 'Mangos', 83, 1.05, 'Local Farm');
