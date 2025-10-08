CREATE OR REPLACE FUNCTION process_order_shipment(
    p_order_id UUID
)
RETURNS VOID AS $$
BEGIN
    -- Move from reserved stock to actual stock reduction
    UPDATE products 
    SET stock = stock - oi.quantity,
        reserved_stock = reserved_stock - oi.quantity,
        updated_at = NOW()
    FROM order_items oi
    WHERE products.id = oi.product_id 
      AND oi.order_id = p_order_id;
    
    -- Update order status to shipped
    UPDATE orders 
    SET status = 'shipped',
        updated_at = NOW()
    WHERE id = p_order_id;
    
    -- Log the outbound transactions
    INSERT INTO transactions (type, product_id, quantity, warehouse_id, reference_number, notes, created_by)
    SELECT 
        'outbound',
        oi.product_id,
        oi.quantity,
        o.warehouse_id,
        o.order_number,
        'Stock shipped for order',
        (SELECT id FROM users WHERE role = 'admin' LIMIT 1)
    FROM order_items oi
    JOIN orders o ON oi.order_id = o.id
    WHERE oi.order_id = p_order_id;
END;
$$ LANGUAGE plpgsql;