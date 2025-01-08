CREATE TABLE payments (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    bank_name VARCHAR(255) NOT NULL,
    bank_account VARCHAR(255) NOT NULL,
    voucher_id UUID NOT NULL REFERENCES vouchers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    date TIMESTAMP,
    invoice VARCHAR(255) NOT NULL,
    payment_option_id UUID NOT NULL REFERENCES payment_options(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP
);