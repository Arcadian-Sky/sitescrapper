CREATE TABLE IF NOT EXISTS products
(
    name String,
    description String,
    price String,
    url String,
    created_at DateTime
) ENGINE = MergeTree()
ORDER BY (created_at); 