CREATE TABLE IF NOT EXISTS "Payment"
(
    "order_uid"     VARCHAR(255) PRIMARY KEY,
    "transaction"   VARCHAR(255) NOT NULL,
    "request_id"    VARCHAR(255),
    "currency"      VARCHAR(255) NOT NULL,
    "provider"      VARCHAR(255) NOT NULL,
    "amount"        INT          NOT NULL,
    "payment_dt"    INT          NOT NULL,
    "bank"          VARCHAR(255) NOT NULL,
    "delivery_cost" INT          NOT NULL,
    "goods_total"   INT          NOT NULL,
    "custom_fee"    INT,
    CONSTRAINT fk_order
        FOREIGN KEY ("order_uid")
            REFERENCES "Order" ("order_uid")
            ON DELETE CASCADE
);
