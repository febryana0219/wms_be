CREATE VIEW products_with_warehouse AS
SELECT 
    p.*,
    w.name as warehouse_name,
    w.code as warehouse_code,
    w.is_active as warehouse_is_active
FROM products p
JOIN warehouses w ON p.warehouse_id = w.id;