CREATE TABLE institutions (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    pendamping_name VARCHAR(255),
    pendamping_phone VARCHAR(50),
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);