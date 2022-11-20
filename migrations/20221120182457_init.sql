-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name CHARACTER VARYING(50) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0 CHECK(balance >= 0)
);

DROP TYPE IF EXISTS transaction_type;
CREATE TYPE transaction_type AS ENUM(
    'write-off',
    'replenishment'
);

DROP TYPE IF EXISTS transaction_status;
CREATE TYPE transaction_status AS ENUM(
    'open',
    'closed',
    'rejected'
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users (id) NOT NULL,
    tr_type transaction_type NOT NULL,
    tr_status transaction_status NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;


