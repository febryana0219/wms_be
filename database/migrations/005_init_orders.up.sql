CREATE TYPE order_status AS ENUM ('pending_payment', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'expired');

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id VARCHAR(100) NOT NULL,
    customer_name VARCHAR(100) NOT NULL,
    status order_status DEFAULT 'pending_payment',
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    notes TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_warehouse_id ON orders(warehouse_id);
CREATE INDEX idx_orders_expires_at ON orders(expires_at);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);

-- Partial index for pending orders that need expiration check
CREATE INDEX idx_orders_pending_expiration ON orders(expires_at) 
WHERE status = 'pending_payment';