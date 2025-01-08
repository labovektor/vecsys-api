CREATE TABLE biodatas (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    participant_id UUID NOT NULL REFERENCES participants(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name VARCHAR(255),
    gender VARCHAR(50),
    phone VARCHAR(50),
    email VARCHAR(255),
    id_number VARCHAR(255),
    id_card_picture VARCHAR(255),
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);