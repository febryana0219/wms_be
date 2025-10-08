CREATE TYPE transaction_type AS ENUM ('transfer', 'pending_payment', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'expired');

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type transaction_type NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    to_warehouse_id UUID REFERENCES warehouses(id),
    reference_number VARCHAR(100),
    notes TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_product_id ON transactions(product_id);
CREATE INDEX idx_transactions_warehouse_id ON transactions(warehouse_id);
CREATE INDEX idx_transactions_to_warehouse_id ON transactions(to_warehouse_id);
CREATE INDEX idx_transactions_created_by ON transactions(created_by);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_reference ON transactions(reference_number);