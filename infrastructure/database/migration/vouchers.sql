CREATE TABLE vouchers (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE,
    voucher VARCHAR(255) NOT NULL,
    desc TEXT NOT NULL,
    seat_available INT DEFAULT 1 NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);