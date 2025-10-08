CREATE OR REPLACE FUNCTION transfer_stock(
    p_product_id UUID,
    p_from_warehouse_id UUID,
    p_to_warehouse_id UUID,
    p_quantity INTEGER,
    p_reference_number VARCHAR(100),
    p_notes TEXT,
    p_created_by UUID
)
RETURNS BOOLEAN AS $$
DECLARE
    source_product_id UUID;
    target_product_id UUID;
BEGIN
    -- Check if source product exists and has enough stock
    SELECT id INTO source_product_id
    FROM products 
    WHERE id = p_product_id 
      AND warehouse_id = p_from_warehouse_id
      AND available_stock >= p_quantity;
    
    IF source_product_id IS NULL THEN
        RETURN FALSE;
    END IF;
    
    -- Check if target product exists in destination warehouse
    SELECT id INTO target_product_id
    FROM products 
    WHERE sku = (SELECT sku FROM products WHERE id = p_product_id)
      AND warehouse_id = p_to_warehouse_id;
    
    -- If target product doesn't exist, create it
    IF target_product_id IS NULL THEN
        INSERT INTO products (name, sku, description, price, stock, warehouse_id, category, min_stock)
        SELECT 
            name, 
            sku, 
            description, 
            price, 
            p_quantity,
            p_to_warehouse_id,
            category,
            min_stock
        FROM products 
        WHERE id = p_product_id
        RETURNING id INTO target_product_id;
    ELSE
        -- Update existing target product stock
        UPDATE products 
        SET stock = stock + p_quantity,
            updated_at = NOW()
        WHERE id = target_product_id;
    END IF;
    
    -- Reduce stock from source warehouse
    UPDATE products 
    SET stock = stock - p_quantity,
        updated_at = NOW()
    WHERE id = source_product_id;
    
    -- Log the transfer transaction
    INSERT INTO transactions (type, product_id, quantity, warehouse_id, to_warehouse_id, reference_number, notes, created_by)
    VALUES ('transfer', p_product_id, p_quantity, p_from_warehouse_id, p_to_warehouse_id, p_reference_number, p_notes, p_created_by);
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;