DELETE FROM products
WHERE sku IN ('LAP-001', 'OFF-001', 'FUR-001');

DELETE FROM categories
WHERE name IN ('Ноутбуки', 'Канцелярия', 'Мебель');

DELETE FROM user_roles
WHERE user_id IN (
    SELECT id FROM users WHERE email = 'test.user@example.com'
);

DELETE FROM users
WHERE email = 'test.user@example.com';

DELETE FROM roles
WHERE code IN ('employee', 'manager', 'admin');