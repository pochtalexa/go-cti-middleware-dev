-- +goose Up

-- создаем таблицу users
CREATE TABLE IF NOT EXISTS users
(
    id        serial PRIMARY KEY,
    login     varchar(40) UNIQUE NOT NULL,
    pass_hash bytea
);

-- +goose Down
-- удаляем таблицу
DROP TABLE IF EXISTS users;
