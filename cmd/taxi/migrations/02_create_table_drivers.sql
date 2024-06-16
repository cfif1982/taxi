-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS drivers(
  id UUID PRIMARY KEY,
  route_id UUID,
	telephone VARCHAR(11) UNIQUE NOT NULL,
  password TEXT,
  name TEXT,
  balance int
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS drivers;
-- +goose StatementEnd
