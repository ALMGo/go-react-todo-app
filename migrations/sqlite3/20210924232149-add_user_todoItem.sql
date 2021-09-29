
-- +migrate Up

CREATE TABLE user (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    username string NOT NULL UNIQUE,
    password string NOT NULL
);

CREATE TABLE todo_item (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER REFERENCES user(id) ON DELETE CASCADE,
  completed INTEGER NOT NULL DEFAULT 0,
  text string NOT NULL,
  created DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  due DATE NOT NULL,
  category string NOT NULL
);

-- +migrate Down

DROP TABLE user;
DROP TABLE  todo_item;
