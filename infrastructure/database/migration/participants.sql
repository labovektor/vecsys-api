CREATE TABLE participants (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE,
    region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name VARCHAR(255) NOT NULL,
    institution_id UUID NOT NULL REFERENCES institutions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    payment_data_id UUID UNIQUE NOT NULL REFERENCES payments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    verified_at TIMESTAMP,
    locked_at TIMESTAMP,
    progress_step participant_progress DEFAULT 'registered' NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);