-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE bioskop (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    lokasi VARCHAR(255) NOT NULL,
    rating FLOAT
);

-- +migrate StatementEnd