
-- +migrate Up

CREATE TABLE user (
    id string NOT NULL UNIQUE PRIMARY KEY,
    email string NOT NULL UNIQUE,
    password string NOT NULL
);

CREATE TABLE todo_item (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  user_id string REFERENCES user(id),
  completed INTEGER NOT NULL,
  text string NOT NULL,
  created DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  due DATE NOT NULL,
  category string NOT NULL
);

-- +migrate Down

DROP TABLE user;
DROP TABLE  todo_item;
