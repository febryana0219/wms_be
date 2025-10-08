-- Password hashes for 'admin123' and 'staff123' (use bcrypt with 12 rounds)
INSERT INTO users (email, password_hash, name, role, warehouse_id, is_active) VALUES
('admin@wms.com', '$2a$10$qZCMQCnZPrHkdsChKTBQXeJS7kV6MMHCeNyBBw2TnWAV0waZ1M66y', 'Admin', 'admin', (SELECT id FROM warehouses WHERE code = 'JKT-MAIN'), true),
('staff@wms.com', '$2a$10$tesHfBohrs/xhVcgrsMq7.E2RNijwftqNTBtcQZ24VgbObcimkZhi', 'Staff', 'staff', (SELECT id FROM warehouses WHERE code = 'JKT-MAIN'), true);