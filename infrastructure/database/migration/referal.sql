CREATE TABLE referals (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE,
    code VARCHAR(255) NOT NULL,
    desc TEXT NOT NULL,
    seat_available INT DEFAULT 1 NOT NULL,
    is_discount BOOLEAN DEFAULT FALSE NOT NULL,
    discount INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);