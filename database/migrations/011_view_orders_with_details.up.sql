CREATE VIEW orders_with_details AS
SELECT 
    o.*,
    w.name as warehouse_name,
    w.code as warehouse_code,
    COUNT(oi.id) as item_count,
    SUM(oi.quantity) as total_quantity
FROM orders o
JOIN warehouses w ON o.warehouse_id = w.id
LEFT JOIN order_items oi ON o.id = oi.order_id
GROUP BY o.id, w.name, w.code;