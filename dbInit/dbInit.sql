-- Perform hard reset

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO public;

COMMENT ON SCHEMA public IS 'standard public schema';

CREATE EXTENSION "uuid-ossp";


-- Set up tables
-- (implicitly using public schema)

CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name character varying(20) UNIQUE NOT NULL,
  last_active timestamptz
);
CREATE INDEX user_name_index ON users (name);

CREATE TABLE messages (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid REFERENCES users,
  timestamp timestamptz NOT NULL DEFAULT now(),
  text text NOT NULL
);
CREATE INDEX message_timestamp_index ON messages (timestamp);

-- End of script
