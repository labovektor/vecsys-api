CREATE TABLE events (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    admin_id UUID NOT NULL REFERENCES admins(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name VARCHAR(255) NOT NULL,
    desc TEXT NOT NULL,
    group_member_num INT DEFAULT 3 NOT NULL,
    icon TEXT,
    participant_target INT NOT NULL,
    period VARCHAR(255) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);