CREATE TABLE IF NOT EXISTS sample
(
    ID VARCHAR PRIMARY KEY NOT NULL,
    BlockchainHeight BIGINT NOT NULL,
    CreatedAt DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS sample_createdat
    ON sample (CreatedAt);