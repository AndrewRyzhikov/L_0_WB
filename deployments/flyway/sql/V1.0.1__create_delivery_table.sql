CREATE TABLE IF NOT EXISTS "Delivery"
(
    "order_uid" VARCHAR(255) PRIMARY KEY,
    "name"      VARCHAR(255) NOT NULL,
    "phone"     VARCHAR(255) NOT NULL,
    "zip"       VARCHAR(255) NOT NULL,
    "city"      VARCHAR(255) NOT NULL,
    "address"   VARCHAR(255) NOT NULL,
    "region"    VARCHAR(255) NOT NULL,
    "email"     VARCHAR(255) NOT NULL,
    CONSTRAINT fk_order
        FOREIGN KEY ("order_uid")
            REFERENCES "Order" ("order_uid")
            ON DELETE CASCADE
);
