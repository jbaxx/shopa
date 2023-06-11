CREATE DATABASE IF NOT EXISTS commercedb;

USE commercedb;

CREATE TABLE IF NOT EXISTS chains (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name string NOT NULL,
    email string NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp()
);

CREATE TABLE IF NOT EXISTS stores (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    chain_id uuid NOT NULL,
    name string NOT NULL,
    phone string NOT NULL,
    address string NOT NULL,
    postal_code string NOT NULL,
    city string NOT NULL,
    country string NOT NULL,
    active bool NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp()
);

CREATE TABLE IF NOT EXISTS inventory (
    store_id uuid NOT NULL,
    item_id uuid NOT NULL,
    quantity int NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp()
);

CREATE TABLE IF NOT EXISTS items (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    sku string NOT NULL,
    name string NOT NULL,
    cost float NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp()
);
