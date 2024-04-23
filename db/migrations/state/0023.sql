-- +migrate Up


-- Add first blob_sequence?
CREATE TABLE state.blob_sequence
(
    index                BIGINT PRIMARY KEY,
    coinbase             VARCHAR,
    final_acc_input_hash VARCHAR,
    first_blob_sequenced  BIGINT,
    last_blob_sequenced   BIGINT,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    block_num          BIGINT NOT NULL REFERENCES state.block (block_num) ON DELETE CASCADE 
);



-- +migrate Down
