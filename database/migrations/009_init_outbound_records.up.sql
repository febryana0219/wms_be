CREATE TYPE outbound_destination_type AS ENUM ('customer', 'warehouse', 'supplier', 'other');

CREATE TABLE outbound_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    destination_type outbound_destination_type DEFAULT 'customer',
    destination_name VARCHAR(100) NOT NULL,
    destination_contact VARCHAR(100),
    reference_number VARCHAR(100),
    unit_price DECIMAL(15,2),
    total_price DECIMAL(15,2),
    notes TEXT,
    shipped_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_outbound_records_product_id ON outbound_records(product_id);
CREATE INDEX idx_outbound_records_warehouse_id ON outbound_records(warehouse_id);
CREATE INDEX idx_outbound_records_destination_type ON outbound_records(destination_type);
CREATE INDEX idx_outbound_records_destination_name ON outbound_records(destination_name);
CREATE INDEX idx_outbound_records_reference_number ON outbound_records(reference_number);
CREATE INDEX idx_outbound_records_shipped_date ON outbound_records(shipped_date DESC);
CREATE INDEX idx_outbound_records_created_by ON outbound_records(created_by);