CREATE VIEW transaction_summary AS
SELECT 
    t.*,
    p.name as product_name,
    p.sku as product_sku,
    w1.name as warehouse_name,
    w2.name as to_warehouse_name,
    u.name as created_by_name,
    u.email as created_by_email
FROM transactions t
JOIN products p ON t.product_id = p.id
JOIN warehouses w1 ON t.warehouse_id = w1.id
LEFT JOIN warehouses w2 ON t.to_warehouse_id = w2.id
JOIN users u ON t.created_by = u.id;