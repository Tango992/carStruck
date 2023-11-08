INSERT INTO brands (name) VALUES ('Toyota'), ('Mitsubishi'), ('Honda');

INSERT INTO categories (name) VALUES ('Hatchback'), ('SUV'), ('MPV');

INSERT INTO catalogs (brand_id, category_id, name, stock, cost)
VALUES 
    (1, 1, 'Yaris', 100, 15),
    (1, 3, 'Avanza', 100, 20),
    (1, 2, 'Fortuner', 100, 30),
    (2, 3, 'Xpander', 100, 23),
    (2, 2, 'Pajero', 100, 32),
    (3, 3, 'BRV', 100, 18),
    (3, 2, 'CRV', 100, 40);