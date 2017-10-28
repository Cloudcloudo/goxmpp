CREATE TABLE users (
  jid         TEXT PRIMARY KEY   NOT NULL,
  password    TEXT               NOT NULL,

  full_name   TEXT               NULL,
  nick_name   TEXT               NULL,
  birthdate   DATE               NULL,
  phone       TEXT               NULL,
  www         TEXT               NULL,
  email       TEXT               NULL,

  company     TEXT               NULL,
  departament TEXT               NULL,
  position    TEXT               NULL,
  role        TEXT               NULL,

  street      TEXT               NULL,
  street_2    TEXT               NULL,
  city        TEXT               NULL,
  state       TEXT               NULL,
  zip_code    TEXT               NULL,
  country     TEXT               NULL,

  about       TEXT               NULL,
  joined      TIMESTAMP DEFAULT now(),
  last_seen   TIMESTAMP DEFAULT now(),

  avatar      TEXT               NULL,

  presence    TEXT               NULL,
  status      TEXT               NULL,

  UNIQUE (jid)
);

CREATE TABLE messages (
  id         TEXT PRIMARY KEY   NOT NULL DEFAULT uuid_generate_v4(),,
  user_jid   TEXT REFERENCES users (jid),
  "from"     TEXT               NOT NULL,
  type       TEXT               NOT NULL,
  subject    TEXT               NULL,
  nick       TEXT               NULL,
  body       TEXT               NOT NULL,
  body_html  TEXT               NULL,

  created_at TIMESTAMP DEFAULT statement_timestamp()
);

CREATE INDEX message_to
  ON messages ("from");

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE contacts (
  id         UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
  user_jid   TEXT REFERENCES users (jid),
  jid        TEXT               NOT NULL,
  "group"    TEXT []            DEFAULT '{}',
  nick       TEXT               NULL,

  subscrbed  TEXT               DEFAULT 'none',

  created_at TIMESTAMP          DEFAULT statement_timestamp()
);

CREATE INDEX contacts_user
  ON contacts (user_jid);


GRANT ALL PRIVILEGES ON DATABASE goxmpp TO goxmpp_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO goxmpp_user;