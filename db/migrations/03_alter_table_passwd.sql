-- +migrate Up
ALTER TABLE users ALTER passwd TYPE VARCHAR(100);
