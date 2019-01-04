-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS tracks (
  user_id TEXT,
  played_at TIMESTAMP,
  duration int,
  track_id TEXT,
  track_name TEXT,
  artist_ids TEXT [],
  artist_names TEXT [],
  PRIMARY KEY(user_id, played_at)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS tracks;
