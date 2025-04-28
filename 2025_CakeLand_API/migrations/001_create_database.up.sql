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
    image_url         VARCHAR(500),
    kg_price          DOUBLE PRECISION                    NOT NULL,
    reviews_count     INT       DEFAULT 0 CHECK (reviews_count >= 0),
    stars_sum         INT       DEFAULT 0 CHECK (stars_sum >= 0),
    description       TEXT                                NOT NULL,
    mass              DOUBLE PRECISION                    NOT NULL,
    discount_kg_price DOUBLE PRECISION CHECK (discount_kg_price >= 0),
    discount_end_time TIMESTAMP,
    date_creation     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_open_for_sale  BOOL      DEFAULT true,
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
    name        VARCHAR(150) UNIQUE NOT NULL,
    image_url   VARCHAR(200)        NOT NULL,
    gender_tags category_gender[] DEFAULT '{}'::category_gender[]
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
    image_url   VARCHAR(200)     NOT NULL,
    content     TEXT             NOT NULL,
    kg_price    DOUBLE PRECISION NOT NULL,
    description TEXT             NOT NULL
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

-- Адреса пользователя
CREATE TABLE IF NOT EXISTS address
(
    id                UUID PRIMARY KEY,                       -- Код адреса
    user_id           UUID             NOT NULL,              -- Ссылка на пользователя, которому принадлежит адрес
    latitude          DOUBLE PRECISION NOT NULL,              -- Географическая широта
    longitude         DOUBLE PRECISION NOT NULL,              -- Географическая долгота
    formatted_address TEXT             NOT NULL,              -- Человеко-читаемый адрес (например, от Google Maps)
    entrance          TEXT,                                   -- Подъезд (необязательное поле)
    floor             TEXT,                                   -- Этаж (необязательное поле)
    apartment         TEXT,                                   -- Квартира (необязательное поле)
    comment           TEXT,                                   -- Комментарий к доставке (например: "оставить у охраны")
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT now(), -- Дата и время создания записи
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT now(), -- Дата и время последнего обновления

    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

-- Статус заказа
CREATE TYPE order_status AS ENUM (
    'pending', -- Ожидает выполнения
    'shipped', -- Отправлен
    'delivered', -- Доставлен
    'cancelled' -- Отменён
    );

-- Перечисление способов оплаты
CREATE TYPE payment_method AS ENUM (
    'cash', -- Наличные
    'io_money' -- ЮMoney
    );

-- Заказ
CREATE TABLE IF NOT EXISTS "order"
(
    id                  uuid PRIMARY KEY,
    total_price         DOUBLE PRECISION CHECK (total_price > 0) NOT NULL,
    delivery_address_id UUID                                     NOT NULL,
    mass                DOUBLE PRECISION                         NOT NULL,
    filling_id          UUID                                     NOT NULL,
    delivery_date       DATE,
    customer_id         UUID                                     NOT NULL,
    seller_id           UUID                                     NOT NULL,
    cake_id             UUID                                     NOT NULL,
    payment_method      payment_method                           NOT NULL DEFAULT 'cash',
    status              order_status                             NOT NULL DEFAULT 'pending',

    FOREIGN KEY (delivery_address_id) REFERENCES "address" (id),
    FOREIGN KEY (cake_id) REFERENCES "cake" (id),
    FOREIGN KEY (customer_id) REFERENCES "user" (id),
    FOREIGN KEY (seller_id) REFERENCES "user" (id),
    FOREIGN KEY (filling_id) REFERENCES "filling" (id)
);

-- Триггеры
CREATE OR REPLACE FUNCTION update_cake_review_stats()
    RETURNS TRIGGER AS
$$
BEGIN
    UPDATE cake
    SET reviews_count = reviews_count + 1,
        stars_sum     = stars_sum + NEW.rating
    WHERE id = NEW.cake_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Функция к триггеру добавления отзыва
CREATE TRIGGER trigger_update_cake_reviews
    AFTER INSERT
    ON feedback
    FOR EACH ROW
EXECUTE FUNCTION update_cake_review_stats();