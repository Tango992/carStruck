INSERT INTO brands (name) VALUES ('Toyota'), ('Mitsubishi'), ('Honda');

INSERT INTO categories (name) VALUES ('Hatchback'), ('SUV'), ('MPV');

INSERT INTO catalogs (brand_id, category_id, name, stock, cost)
VALUES 
    (1, 1, 'Yaris', 100, 600000),
    (1, 3, 'Avanza', 100, 500000),
    (1, 2, 'Fortuner', 100, 1000000),
    (2, 3, 'Xpander', 100, 700000),
    (2, 2, 'Pajero', 100, 1250000),
    (3, 3, 'Brv', 100, 550000),
    (3, 2, 'Crv', 100, 800000);