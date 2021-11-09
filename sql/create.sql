CREATE TYPE role AS ENUM ('member', 'admin', 'owner');

CREATE TABLE users (
  user_id  SERIAL PRIMARY KEY,
  role     role   NOT NULL,
  name     TEXT   NOT NULL UNIQUE,
  password TEXT   NOT NULL
);

CREATE TABLE threads (
  thread_id SERIAL PRIMARY KEY,
  title     TEXT   NOT NULL
);

CREATE TABLE comments (
  comment_id SERIAL,
  thread_id  INTEGER,
  content    TEXT NOT NULL,
  PRIMARY KEY (comment_id, thread_id),
  FOREIGN KEY (thread_id) REFERENCES threads(thread_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE categories (
  category_id SERIAL PRIMARY KEY,
  name        TEXT   NOT NULL UNIQUE
);

CREATE TABLE tags (
  tag_id SERIAL PRIMARY KEY,
  name   TEXT   NOT NULL UNIQUE
);

CREATE TABLE post_threads (
  thread_id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP,
  FOREIGN KEY (thread_id) REFERENCES threads(thread_id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE post_comments (
  comment_id INTEGER, 
  thread_id INTEGER,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP,
  PRIMARY KEY (comment_id, thread_id),
  FOREIGN KEY (comment_id, thread_id) REFERENCES comments(comment_id, thread_id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE link_categories (
  thread_id INTEGER PRIMARY KEY,
  category_id INTEGER NOT NULL,
  FOREIGN KEY (thread_id) REFERENCES threads(thread_id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE add_tags (
  thread_id INTEGER,
  tag_id INTEGER,
  PRIMARY KEY (thread_id, tag_id),
  FOREIGN KEY (thread_id) REFERENCES threads(thread_id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE ON UPDATE CASCADE
);
