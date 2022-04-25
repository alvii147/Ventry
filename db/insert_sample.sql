\connect ventrydb;

INSERT INTO Inventory (title, quantity, price, owner, supplier, shipper)
VALUES
    ('Cream Cheese', 15, 3.69, 'Walmart', 'PHILADELPHIA', 'Instacart'),
    ('Eggs', 8, 4.29, 'Walmart', 'Local Farms', 'Instacart'),
    ('Yogurt', 3, 5.52, 'No Frills', 'Astro', 'Uber')
