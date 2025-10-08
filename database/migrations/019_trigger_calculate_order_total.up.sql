CREATE OR REPLACE FUNCTION calculate_order_total()
RETURNS TRIGGER AS $$
BEGIN
    -- Update order total when order items change
    UPDATE orders 
    SET total_amount = (
        SELECT COALESCE(SUM(total_price), 0)
        FROM order_items 
        WHERE order_id = COALESCE(NEW.order_id, OLD.order_id)
    )
    WHERE id = COALESCE(NEW.order_id, OLD.order_id);
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_calculate_order_total_insert
    AFTER INSERT ON order_items
    FOR EACH ROW
    EXECUTE FUNCTION calculate_order_total();

CREATE TRIGGER trigger_calculate_order_total_update
    AFTER UPDATE ON order_items
    FOR EACH ROW
    EXECUTE FUNCTION calculate_order_total();

CREATE TRIGGER trigger_calculate_order_total_delete
    AFTER DELETE ON order_items
    FOR EACH ROW
    EXECUTE FUNCTION calculate_order_total();