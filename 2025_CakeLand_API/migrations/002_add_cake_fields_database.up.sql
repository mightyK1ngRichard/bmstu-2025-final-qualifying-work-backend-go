ALTER TABLE "cake"
    ADD COLUMN discount_kg_price DOUBLE PRECISION CHECK (discount_kg_price >= 0),
    ADD COLUMN discount_end_time TIMESTAMP,
    ADD COLUMN date_creation     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL;