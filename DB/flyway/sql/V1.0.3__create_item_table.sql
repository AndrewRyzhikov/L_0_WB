CREATE TABLE IF NOT EXISTS "Item"
(
    "id"           SERIAL PRIMARY KEY,
    "order_uid"    VARCHAR(255) NOT NULL,
    "chrt_id"      INT          NOT NULL,
    "track_number" VARCHAR(255) NOT NULL,
    "price"        INT          NOT NULL,
    "rid"          VARCHAR(255) NOT NULL,
    "name"         VARCHAR(255) NOT NULL,
    "sale"         INT          NOT NULL,
    "size"         VARCHAR(255) NOT NULL,
    "total_price"  INT          NOT NULL,
    "nm_id"        INT          NOT NULL,
    "brand"        VARCHAR(255) NOT NULL,
    "status"       INT          NOT NULL,
    CONSTRAINT fk_order
        FOREIGN KEY ("order_uid")
            REFERENCES "Order" ("order_uid")
            ON DELETE CASCADE
);
