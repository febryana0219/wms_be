CREATE TABLE inbound_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    supplier_name VARCHAR(100) NOT NULL,
    supplier_contact VARCHAR(100),
    reference_number VARCHAR(100),
    unit_cost DECIMAL(15,2),
    total_cost DECIMAL(15,2),
    notes TEXT,
    received_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_inbound_records_product_id ON inbound_records(product_id);
CREATE INDEX idx_inbound_records_warehouse_id ON inbound_records(warehouse_id);
CREATE INDEX idx_inbound_records_supplier_name ON inbound_records(supplier_name);
CREATE INDEX idx_inbound_records_reference_number ON inbound_records(reference_number);
CREATE INDEX idx_inbound_records_received_date ON inbound_records(received_date DESC);
CREATE INDEX idx_inbound_records_created_by ON inbound_records(created_by);