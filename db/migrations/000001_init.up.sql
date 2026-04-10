CREATE TABLE IF NOT EXISTS users(
  id uuid PRIMARY KEY DEFAULT uuidv7(),
  name TEXT NOT NULL,
  hash bytea NOT NULL
);

CREATE TABLE IF NOT EXISTS todos(
  id uuid PRIMARY KEY DEFAULT uuidv7(),
  text TEXT NOT NULL,
  wip_at timestamptz,
  completed_at timestamptz,
  user_id uuid NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_wip_at ON todos(wip_at);
CREATE INDEX IF NOT EXISTS idx_completed_at ON todos(completed_at);
CREATE INDEX IF NOT EXISTS idx_user_id ON todos(user_id);
