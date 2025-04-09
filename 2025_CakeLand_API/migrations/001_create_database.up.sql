-- Пользователь
CREATE TABLE IF NOT EXISTS "user"
(
    id                 UUID PRIMARY KEY,
    fio                VARCHAR(150),
    address            TEXT,
    nickname           VARCHAR(50) UNIQUE  NOT NULL,
    image_url          VARCHAR(200),
    header_image_url   VARCHAR(200),
    mail               VARCHAR(200) UNIQUE NOT NULL,
    password_hash      VARCHAR(100) UNIQUE NOT NULL,
    phone              VARCHAR(11),
    card_number        VARCHAR(16),
    refresh_tokens_map JSONB
);

-- Торт
CREATE TABLE IF NOT EXISTS cake
(
    id                UUID PRIMARY KEY,
    name              VARCHAR(150)                        NOT NULL,
    image_url         VARCHAR(200),
    kg_price          DOUBLE PRECISION                    NOT NULL,
    rating            INT       DEFAULT 0 CHECK (rating >= 0 AND rating <= 5),
    description       TEXT,
    mass              DOUBLE PRECISION,
    discount_kg_price DOUBLE PRECISION CHECK (discount_kg_price >= 0),
    discount_end_time TIMESTAMP,
    date_creation     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_open_for_sale  BOOL,
    owner_id          UUID                                NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES "user" (id)
);

-- Фотографии торта
CREATE TABLE IF NOT EXISTS "cake_images"
(
    id        UUID PRIMARY KEY,
    image_url VARCHAR(200),
    cake_id   UUID,
    FOREIGN KEY (cake_id) REFERENCES "cake" (id)
);

-- Отзыв
CREATE TABLE IF NOT EXISTS feedback
(
    id            UUID PRIMARY KEY,
    text          TEXT,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    rating        INT       DEFAULT 0 CHECK (rating >= 0 AND rating <= 5),
    cake_id       UUID NOT NULL,
    author_id     UUID NOT NULL,
    FOREIGN KEY (cake_id) REFERENCES "cake" (id),
    FOREIGN KEY (author_id) REFERENCES "user" (id)
);

-- Пол категории
CREATE TYPE category_gender AS ENUM ('male', 'female', 'child');

-- Категория
CREATE TABLE IF NOT EXISTS category
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(150) NOT NULL,
    image_url   VARCHAR(200),
    gender_tags category_gender[]
);

-- Категории торта (М-М)
CREATE TABLE IF NOT EXISTS cake_category
(
    id            UUID PRIMARY KEY,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category_id   UUID NOT NULL,
    cake_id       UUID NOT NULL,
    FOREIGN KEY (category_id) REFERENCES "category" (id),
    FOREIGN KEY (cake_id) REFERENCES "cake" (id)
);

-- Начинка
CREATE TABLE IF NOT EXISTS filling
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(150)     NOT NULL,
    image_url   VARCHAR(200),
    content     TEXT,
    kg_price    DOUBLE PRECISION NOT NULL,
    description TEXT
);

-- Начинки торта (М-М)
CREATE TABLE IF NOT EXISTS cake_filling
(
    id         uuid PRIMARY KEY,
    cake_id    UUID NOT NULL,
    filling_id UUID NOT NULL,
    FOREIGN KEY (cake_id) REFERENCES "cake" (id),
    FOREIGN KEY (filling_id) REFERENCES "filling" (id)
);

-- Сообщение
CREATE TABLE IF NOT EXISTS message
(
    id            uuid PRIMARY KEY,
    text          TEXT,
    date_creation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    owner_id      UUID NOT NULL,
    receiver_id   UUID NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES "user" (id),
    FOREIGN KEY (receiver_id) REFERENCES "user" (id)
);

-- Статус заказа
CREATE TYPE order_status AS ENUM (
    'pending', -- Ожидает выполнения
    'shipped', -- Отправлен
    'delivered', -- Доставлен
    'cancelled' -- Отменён
    );

-- Заказ
CREATE TABLE IF NOT EXISTS "order"
(
    id               uuid PRIMARY KEY,
    price            DOUBLE PRECISION CHECK (price > 0),
    delivery_address TEXT,
    delivery_date    DATE,
    customer_id      UUID         NOT NULL,
    seller_id        UUID         NOT NULL,
    cake_id          UUID         NOT NULL,
    status           order_status NOT NULL DEFAULT 'pending',
    FOREIGN KEY (cake_id) REFERENCES "cake" (id),
    FOREIGN KEY (customer_id) REFERENCES "user" (id),
    FOREIGN KEY (seller_id) REFERENCES "user" (id)
);
