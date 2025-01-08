CREATE TABLE payment_options (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE,
    provider VARCHAR(255) NOT NULL,
    account VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    as_qr BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);