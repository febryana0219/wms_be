CREATE OR REPLACE FUNCTION release_stock_for_order(
    p_order_id UUID
)
RETURNS VOID AS $$
BEGIN
    -- Release reserved stock for all items in the order
    UPDATE products 
    SET reserved_stock = reserved_stock - oi.quantity,
        updated_at = NOW()
    FROM order_items oi
    WHERE products.id = oi.product_id 
      AND oi.order_id = p_order_id;
    
    -- Log the release transactions
    INSERT INTO transactions (type, product_id, quantity, warehouse_id, reference_number, notes, created_by)
    SELECT 
        'release',
        oi.product_id,
        oi.quantity,
        o.warehouse_id,
        o.order_number,
        'Stock released due to order cancellation/expiration',
        (SELECT id FROM users WHERE role = 'admin' LIMIT 1)
    FROM order_items oi
    JOIN orders o ON oi.order_id = o.id
    WHERE oi.order_id = p_order_id;
END;
$$ LANGUAGE plpgsql;