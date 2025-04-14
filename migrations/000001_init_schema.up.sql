CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pickup_points (
    id SERIAL PRIMARY KEY,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW(),
    city VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS receipts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    pickup_point_id INTEGER NOT NULL REFERENCES pickup_points(id),
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress',
    CONSTRAINT status_check CHECK (status IN ('in_progress', 'closed'))
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    received_at TIMESTAMP NOT NULL DEFAULT NOW(),
    type VARCHAR(50) NOT NULL,
    receipt_id INTEGER NOT NULL REFERENCES receipts(id),
    addition_order INTEGER NOT NULL,
    CONSTRAINT type_check CHECK (type IN ('электроника', 'одежда', 'обувь'))
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_receipts_pickup_point_id ON receipts(pickup_point_id);
CREATE INDEX IF NOT EXISTS idx_products_receipt_id ON products(receipt_id);
CREATE INDEX IF NOT EXISTS idx_receipts_created_at ON receipts(created_at); 