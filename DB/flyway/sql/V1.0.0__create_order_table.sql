CREATE TABLE IF NOT EXISTS "Order"
(
    "order_uid"          VARCHAR(255) PRIMARY KEY,
    "track_number"       VARCHAR(255) NOT NULL,
    "entry"              VARCHAR(255) NOT NULL,
    "locale"             VARCHAR(255) NOT NULL,
    "internal_signature" VARCHAR(255),
    "customer_id"        VARCHAR(255) NOT NULL,
    "delivery_service"   VARCHAR(255) NOT NULL,
    "shardkey"           VARCHAR(255) NOT NULL,
    "sm_id"              INT          NOT NULL,
    "date_created"       TIMESTAMP    NOT NULL,
    "oof_shard"          VARCHAR(255) NOT NULL
);