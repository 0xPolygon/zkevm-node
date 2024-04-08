-- +migrate Up
ALTER TABLE state.block
    ADD COLUMN IF NOT EXISTS checked BOOL NOT NULL DEFAULT TRUE;

-- +migrate Down
ALTER TABLE state.receipt
    DROP COLUMN IF EXISTS checked;