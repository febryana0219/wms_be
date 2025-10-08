CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(100),
    description TEXT,
    price DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    min_stock INTEGER NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    reserved_stock INTEGER NOT NULL DEFAULT 0,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Computed column for available stock
ALTER TABLE products ADD COLUMN available_stock INTEGER 
GENERATED ALWAYS AS (stock - reserved_stock) STORED;

-- Indexes
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_warehouse_id ON products(warehouse_id);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_stock ON products(stock);
CREATE INDEX idx_products_available_stock ON products(available_stock);
CREATE INDEX idx_products_is_active ON products(is_active);

-- Full text search index for name and description
CREATE INDEX idx_products_search ON products USING gin(to_tsvector('english', name || ' ' || COALESCE(description, '')));