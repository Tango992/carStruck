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

-- For testing samples
INSERT INTO brands (name) VALUES ('Toyota');

INSERT INTO categories (name) VALUES ('MPV');

INSERT INTO catalogs (brand_id, category_id, name, stock, cost)
VALUES (1, 1, 'Avanza', 100, 500000);

INSERT INTO users (full_name, email, password, birth, address)
VALUES ('John Doe', 'john@mail.com', '$2a$10$.F.xwVpESnzJo4gFvzhVM.7WBd7szXRn/BvlRoHISLyshSh32WEUu', '2001-01-01', 'Perkantoran The Breeze BSD');

INSERT INTO verifications (user_id, token)
VALUES (1, 'abcd');

INSERT INTO orders (user_id, catalog_id, rent_date, return_date)
VALUES (1, 1, '2023-11-10', '2023-11-20');

INSERT INTO payments (order_id, invoice_id, amount, invoice_url, status, created_at)
VALUES (1, 'abcd', 5000000, 'https://google.com', 'PENDING', '2023-11-01 00:00:00');