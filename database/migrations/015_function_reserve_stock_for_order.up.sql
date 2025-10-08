CREATE OR REPLACE FUNCTION reserve_stock_for_order(
    p_order_id UUID,
    p_product_id UUID,
    p_quantity INTEGER
)
RETURNS BOOLEAN AS $$
BEGIN
    -- Check if enough stock is available
    IF (SELECT available_stock FROM products WHERE id = p_product_id) < p_quantity THEN
        RETURN FALSE;
    END IF;
    
    -- Reserve the stock
    UPDATE products 
    SET reserved_stock = reserved_stock + p_quantity,
        updated_at = NOW()
    WHERE id = p_product_id;
    
    -- Log the reservation transaction
    INSERT INTO transactions (type, product_id, quantity, warehouse_id, reference_number, notes, created_by)
    SELECT 
        'checkout',
        p_product_id,
        p_quantity,
        o.warehouse_id,
        o.order_number,
        'Stock reserved for order',
        (SELECT id FROM users WHERE role = 'admin' LIMIT 1)
    FROM orders o 
    WHERE o.id = p_order_id;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;