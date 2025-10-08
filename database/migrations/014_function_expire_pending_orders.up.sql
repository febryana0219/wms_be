CREATE OR REPLACE FUNCTION expire_pending_orders()
RETURNS INTEGER AS $$
DECLARE
    expired_count INTEGER;
BEGIN
    -- Update expired orders and release reserved stock
    WITH expired_orders AS (
        UPDATE orders 
        SET status = 'expired',
            updated_at = NOW()
        WHERE status = 'pending_payment' 
          AND expires_at < NOW()
        RETURNING id
    ),
    order_items_to_release AS (
        SELECT oi.product_id, SUM(oi.quantity) as total_quantity
        FROM order_items oi
        JOIN expired_orders eo ON oi.order_id = eo.id
        GROUP BY oi.product_id
    )
    UPDATE products 
    SET reserved_stock = reserved_stock - oitr.total_quantity,
        updated_at = NOW()
    FROM order_items_to_release oitr
    WHERE products.id = oitr.product_id;
    
    -- Get count of expired orders
    GET DIAGNOSTICS expired_count = ROW_COUNT;
    
    -- Log transactions for released stock
    INSERT INTO transactions (type, product_id, quantity, warehouse_id, reference_number, notes, created_by)
    SELECT 
        'release',
        oi.product_id,
        oi.quantity,
        o.warehouse_id,
        o.order_number,
        'Auto-released due to order expiration',
        (SELECT id FROM users WHERE role = 'admin' LIMIT 1)
    FROM order_items oi
    JOIN orders o ON oi.order_id = o.id
    WHERE o.status = 'expired'
      AND o.updated_at >= NOW() - INTERVAL '1 minute'; -- Only just expired orders
    
    RETURN expired_count;
END;
$$ LANGUAGE plpgsql;