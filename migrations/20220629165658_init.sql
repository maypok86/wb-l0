-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS deliveries(
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar NOT NULL,
    phone varchar NOT NULL,
    zip varchar NOT NULL,
    city varchar NOT NULL,
    address varchar NOT NULL,
    region varchar NOT NULL,
    email varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS payments(
    transaction varchar PRIMARY KEY,
    request_id varchar NOT NULL,
    currency varchar NOT NULL,
    provider varchar NOT NULL,
    amount int NOT NULL,
    payment_dt int NOT NULL,
    bank varchar NOT NULL,
    delivery_cost int NOT NULL,
    goods_total int NOT NULL,
    custom_fee int NOT NULL
);

CREATE TABLE IF NOT EXISTS orders(
    order_uid varchar PRIMARY KEY,
    track_number varchar NOT NULL UNIQUE,
    entry varchar NOT NULL,
    delivery_id int NOT NULL,
    locale varchar NOT NULL,
    internal_signature varchar,
    customer_id varchar NOT NULL,
    delivery_service varchar NOT NULL,
    shardkey varchar NOT NULL,
    sm_id int NOT NULL,
    date_created timestamptz NOT NULL DEFAULT now(),
    oof_shard varchar NOT NULL,

    FOREIGN KEY (delivery_id) REFERENCES deliveries(id),
    FOREIGN KEY (order_uid) REFERENCES payments(transaction)
);

CREATE TABLE IF NOT EXISTS items(
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chrt_id int NOT NULL,
    track_number varchar NOT NULL,
    price int NOT NULL,
    rid varchar NOT NULL,
    name varchar NOT NULL,
    sale int NOT NULL,
    size varchar NOT NULL,
    total_price int NOT NULL,
    nm_id int NOT NULL,
    brand varchar NOT NULL,
    status int NOT NULL,

    FOREIGN KEY (track_number) REFERENCES orders(track_number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
