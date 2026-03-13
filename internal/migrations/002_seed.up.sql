INSERT INTO roles (code, name)
VALUES
    ('employee', 'Сотрудник'),
    ('manager', 'Менеджер'),
    ('admin', 'Администратор');

INSERT INTO users (email, full_name)
VALUES
    ('test.user@example.com', 'Тестовый пользователь');

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email = 'test.user@example.com'
  AND r.code = 'employee';

INSERT INTO categories (name, description)
VALUES
    ('Ноутбуки', 'Ноутбуки и аксессуары'),
    ('Канцелярия', 'Офисные принадлежности'),
    ('Мебель', 'Офисная мебель');

INSERT INTO products (category_id, name, description, sku, unit)
SELECT c.id, 'Lenovo ThinkPad E14', 'Ноутбук для работы', 'LAP-001', 'pcs'
FROM categories c
WHERE c.name = 'Ноутбуки';

INSERT INTO products (category_id, name, description, sku, unit)
SELECT c.id, 'Шариковая ручка', 'Синяя ручка', 'OFF-001', 'pcs'
FROM categories c
WHERE c.name = 'Канцелярия';

INSERT INTO products (category_id, name, description, sku, unit)
SELECT c.id, 'Офисное кресло', 'Кресло для рабочего места', 'FUR-001', 'pcs'
FROM categories c
WHERE c.name = 'Мебель';