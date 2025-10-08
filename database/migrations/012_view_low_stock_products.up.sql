CREATE VIEW low_stock_products AS
SELECT 
    p.*,
    w.name as warehouse_name,
    (p.min_stock - p.available_stock) as shortage_quantity
FROM products p
JOIN warehouses w ON p.warehouse_id = w.id
WHERE p.available_stock <= p.min_stock
  AND p.is_active = true
  AND w.is_active = true
ORDER BY p.available_stock ASC;