-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS routes(
	id int PRIMARY KEY,
	name VARCHAR(7) UNIQUE NOT NULL,
  points text
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS routes;
-- +goose StatementEnd
