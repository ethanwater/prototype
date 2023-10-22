DROP TABLE IF EXISTS users;
CREATE TABLE users(
  id         INT AUTO_INCREMENT NOT NULL,
  name      VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO users
  (name)
VALUES
  ('Lexi Gray'),
  ('Jeru Steps');