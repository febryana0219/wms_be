INSERT INTO products (name, sku, description, price, stock, reserved_stock, min_stock, warehouse_id, category, is_active) VALUES
('Vitamin C Serum Brightening', 'VITC-SER-001', 'Anti-aging serum with 20% Vitamin C for brightening skin', 350000.00, 125, 12, 20, 
 (SELECT id FROM warehouses WHERE code = 'JKT-MAIN'), 'Serum', true),
('Hyaluronic Acid Moisturizer', 'HYAL-MOIST-001', 'Intensive hydrating moisturizer with hyaluronic acid', 275000.00, 85, 8, 15, 
 (SELECT id FROM warehouses WHERE code = 'JKT-MAIN'), 'Moisturizer', true),
('Niacinamide Toner 10%', 'NIAC-TON-001', 'Gentle toner with 10% niacinamide for pore refinement', 195000.00, 156, 0, 25, 
 (SELECT id FROM warehouses WHERE code = 'JKT-MAIN'), 'Toner', true),
('Gentle Foaming Cleanser', 'FOAM-CLEAN-001', 'pH-balanced foaming cleanser for all skin types', 165000.00, 234, 5, 30, 
 (SELECT id FROM warehouses WHERE code = 'JKT-MAIN'), 'Cleanser', true),
('Anti-Aging Night Cream', 'NIGHT-CREAM-001', 'Rich night cream with retinol and peptides', 425000.00, 8, 0, 15, 
 (SELECT id FROM warehouses WHERE code = 'JKT-MAIN'), 'Moisturizer', true);