-- +goose Up
-- +goose StatementBegin
CREATE TABLE ADDRESSES (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(50),
    is_default BOOLEAN DEFAULT FALSE,
    line_1 VARCHAR(255) NOT NULL,
    line_2 VARCHAR(255),
    postal_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_adresses_user_id ON addresses(user_id);

CREATE TRIGGER trigger_update_addresses_updated_at
BEFORE UPDATE ON addresses
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_adresses_user_id;
DROP TRIGGER IF EXISTS trigger_update_addresses_updated_at ON addresses;
DROP TABLE IF EXISTS addresses;
-- +goose StatementEnd
