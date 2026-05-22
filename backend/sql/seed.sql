INSERT INTO juice (id, name, description, price, created_at, updated_at, stock) VALUES
(gen_random_uuid(), 'Apple Juice', 'Fresh pressed apple juice', 399, NOW(), NOW(), 50),
(gen_random_uuid(), 'Orange Juice', 'Squeezed oranges', 499, NOW(), NOW(), 50);