CREATE TABLE IF NOT EXISTS "cake_images"
(
    id        UUID PRIMARY KEY,
    image_url VARCHAR(200),
    cake_id   UUID,
    FOREIGN KEY (cake_id) REFERENCES "cake" (id)
);